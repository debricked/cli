package nuget

import (
	"fmt"
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

	sanitized := util.NewOrderedObject()
	sanitized.Set("dependencies", sanitizeFrameworkGroups(dependencies, renames))

	return util.EncodeJSON(sanitized)
}

func findProjectRenames(dependencies *util.OrderedObject, context *util.FileContext) (map[string]string, error) {
	renames := map[string]string{}
	used := allPackageNames(dependencies)
	nextPlaceholder := 1
	nextProjectPlaceholder := func() string {
		for {
			placeholder := fmt.Sprintf("redacted-project-%d", nextPlaceholder)
			nextPlaceholder++
			lower := strings.ToLower(placeholder)
			if !used[lower] {
				used[lower] = true

				return placeholder
			}
		}
	}

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
			if !isType || !strings.EqualFold(entryType, "Project") {
				continue
			}

			lookup := strings.ToLower(packageName)
			if _, exists := renames[lookup]; !exists {
				placeholder := nextProjectPlaceholder()
				renames[lookup] = placeholder
				context.RecordPlaceholder(placeholder, packageName)
			}
		}
	}

	return renames, nil
}

func allPackageNames(dependencies *util.OrderedObject) map[string]bool {
	used := map[string]bool{}
	for _, framework := range dependencies.Keys {
		group, isObject := dependencies.Values[framework].(*util.OrderedObject)
		if !isObject {
			continue
		}
		for _, packageName := range group.Keys {
			used[strings.ToLower(packageName)] = true
		}
	}

	return used
}

func sanitizeFrameworkGroups(dependencies *util.OrderedObject, renames map[string]string) *util.OrderedObject {
	sanitized := util.NewOrderedObject()
	for _, framework := range dependencies.Keys {
		group := dependencies.Values[framework].(*util.OrderedObject)
		sanitized.Set(framework, sanitizePackageEntries(group, renames))
	}

	return sanitized
}

func sanitizePackageEntries(group *util.OrderedObject, renames map[string]string) *util.OrderedObject {
	sanitized := util.NewOrderedObject()
	for _, packageName := range group.Keys {
		entry := group.Values[packageName].(*util.OrderedObject)
		sanitized.Set(renameIfProject(packageName, renames), sanitizePackageEntry(entry, renames))
	}

	return sanitized
}

func sanitizePackageEntry(entry *util.OrderedObject, renames map[string]string) *util.OrderedObject {
	sanitized := util.NewOrderedObject()
	for _, key := range []string{"type", "requested", "resolved", "contentHash", "dependencies"} {
		if _, exists := entry.Values[key]; !exists {
			continue
		}
		switch key {
		case "type", "requested", "resolved", "contentHash":
			sanitized.Set(key, entry.Values[key])
		case "dependencies":
			dependencies, isObject := entry.Values[key].(*util.OrderedObject)
			if isObject {
				sanitized.Set(key, sanitizeDependencyVersions(dependencies, renames))
			}
		}
	}

	return sanitized
}

func sanitizeDependencyVersions(dependencies *util.OrderedObject, renames map[string]string) *util.OrderedObject {
	sanitized := util.NewOrderedObject()
	for _, packageName := range dependencies.Keys {
		sanitized.Set(renameIfProject(packageName, renames), dependencies.Values[packageName])
	}

	return sanitized
}

func renameIfProject(packageName string, renames map[string]string) string {
	if renamed, found := renames[strings.ToLower(packageName)]; found {
		return renamed
	}

	return packageName
}
