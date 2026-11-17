package nuget

import (
	"strings"

	"github.com/debricked/cli/internal/sanitize/pm/util"
)

// sanitizePackagesLockJSON keeps only the dependency graph fields used by the backend parser.
func sanitizePackagesLockJSON(content []byte, context *util.FileContext) ([]byte, error) {
	lockFile, err := util.ParseJSONObject(content)
	if err != nil {
		return nil, util.ErrInvalidFormat
	}

	dependencies, isObject := lockFile.Values["dependencies"].(*util.OrderedObject)
	if !isObject {
		return nil, util.ErrInvalidFormat
	}

	renames, err := findProjectRenames(dependencies, context)
	if err != nil {
		return nil, err
	}

	groups, err := sanitizeFrameworkGroups(dependencies, renames)
	if err != nil {
		return nil, err
	}
	sanitized := util.NewOrderedObject()
	sanitized.Set("dependencies", groups)

	return util.EncodeJSON(sanitized)
}

// findProjectRenames gives each `type: Project` entry, the customer's own project, a placeholder
// that is unique across the scan. The backend reports these as packages, and /sbom/merge matches
// on purl, so a placeholder shared with another project would merge the two. Names are matched
// ignoring case, as NuGet does.
func findProjectRenames(dependencies *util.OrderedObject, context *util.FileContext) (map[string]string, error) {
	renames := map[string]string{}
	for _, framework := range dependencies.Keys {
		group, isObject := dependencies.Values[framework].(*util.OrderedObject)
		if !isObject {
			return nil, util.ErrInvalidFormat
		}
		for _, packageName := range group.Keys {
			entry, isEntry := group.Values[packageName].(*util.OrderedObject)
			if !isEntry {
				return nil, util.ErrInvalidFormat
			}
			entryType, isType := entry.Values["type"].(string)
			if !isType {
				return nil, util.ErrInvalidFormat
			}
			lookup := strings.ToLower(packageName)
			if _, exists := renames[lookup]; strings.EqualFold(entryType, "Project") && !exists {
				renames[lookup] = context.ProjectPlaceholder(packageName)
			}
		}
	}

	return renames, nil
}

func sanitizeFrameworkGroups(dependencies *util.OrderedObject, renames map[string]string) (*util.OrderedObject, error) {
	sanitized := util.NewOrderedObject()
	for _, framework := range dependencies.Keys {
		group, isObject := dependencies.Values[framework].(*util.OrderedObject)
		if !isObject {
			return nil, util.ErrInvalidFormat
		}
		entries, err := sanitizePackageEntries(group, renames)
		if err != nil {
			return nil, err
		}
		sanitized.Set(framework, entries)
	}

	return sanitized, nil
}

func sanitizePackageEntries(group *util.OrderedObject, renames map[string]string) (*util.OrderedObject, error) {
	sanitized := util.NewOrderedObject()
	for _, packageName := range group.Keys {
		entry, isObject := group.Values[packageName].(*util.OrderedObject)
		if !isObject {
			return nil, util.ErrInvalidFormat
		}
		sanitizedEntry, err := sanitizePackageEntry(entry, renames)
		if err != nil {
			return nil, err
		}
		sanitized.Set(renameIfProject(packageName, renames), sanitizedEntry)
	}

	return sanitized, nil
}

// packageEntryStrings are the keys of an entry that the backend reads, all of them strings.
// contentHash has to stay: the parser finds each package's line by searching for it.
var packageEntryStrings = []string{"type", "requested", "resolved", "contentHash"}

// sanitizePackageEntry keeps only string values, so an unexpected nested value can't be copied.
func sanitizePackageEntry(entry *util.OrderedObject, renames map[string]string) (*util.OrderedObject, error) {
	sanitized := util.NewOrderedObject()
	for _, key := range packageEntryStrings {
		value, exists := entry.Values[key]
		if !exists {
			continue
		}
		text, isString := value.(string)
		if !isString {
			return nil, util.ErrInvalidFormat
		}
		sanitized.Set(key, text)
	}

	value, exists := entry.Values["dependencies"]
	if !exists {
		return sanitized, nil
	}
	dependencies, isObject := value.(*util.OrderedObject)
	if !isObject {
		return nil, util.ErrInvalidFormat
	}
	sanitizedDependencies := util.NewOrderedObject()
	for _, name := range dependencies.Keys {
		version, isString := dependencies.Values[name].(string)
		if !isString {
			return nil, util.ErrInvalidFormat
		}
		sanitizedDependencies.Set(renameIfProject(name, renames), version)
	}
	sanitized.Set("dependencies", sanitizedDependencies)

	return sanitized, nil
}

func renameIfProject(packageName string, renames map[string]string) string {
	if renamed, found := renames[strings.ToLower(packageName)]; found {
		return renamed
	}

	return packageName
}
