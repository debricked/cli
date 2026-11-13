package npm

import (
	"strings"

	"github.com/debricked/cli/internal/sanitize/pm/util"
)

const (
	packageLockPathFormat    = "./redacted-%d"
	packageLockModulesPrefix = "node_modules/"
	packageLockModulesInfix  = "/node_modules/"
)

// packageLockSanitizer holds the path placeholders of one file, so that an entry key and the
// resolved value of a link pointing at it get the same placeholder.
type packageLockSanitizer struct {
	context *util.FileContext
	paths   map[string]string
}

// sanitizePackageLock handles lockfile versions 1 to 3. Version 1 keeps its packages in nested
// `dependencies`, versions 2 and 3 in `packages`, and version 2 has both.
func sanitizePackageLock(content []byte, context *util.FileContext) ([]byte, error) {
	lock, err := util.ParseJSONObject(content)
	if err != nil {
		return nil, util.ErrInvalidFormat
	}
	sanitizer := &packageLockSanitizer{context: context, paths: map[string]string{}}

	sanitized := util.NewOrderedObject()
	recognised := false
	for _, key := range lock.Keys {
		value := lock.Values[key]
		switch key {
		case "lockfileVersion":
			recognised = true
			sanitized.Set(key, value)
		case "packages":
			packages, isObject := value.(*util.OrderedObject)
			if !isObject {
				return nil, util.ErrInvalidFormat
			}
			recognised = true
			sanitized.Set(key, sanitizer.sanitizePackages(packages))
		case "dependencies":
			dependencies, isObject := value.(*util.OrderedObject)
			if !isObject {
				return nil, util.ErrInvalidFormat
			}
			recognised = true
			sanitized.Set(key, sanitizer.sanitizeList(dependencies))
		}
	}
	if !recognised {
		return nil, util.ErrInvalidFormat
	}

	return util.EncodeJSON(sanitized)
}

func (s *packageLockSanitizer) sanitizePackages(packages *util.OrderedObject) *util.OrderedObject {
	// Keys get their numbers first, so numbering doesn't depend on where a link appears.
	for _, key := range packages.Keys {
		s.maskPath(key)
	}

	sanitized := util.NewOrderedObject()
	for _, key := range packages.Keys {
		entry, isObject := packages.Values[key].(*util.OrderedObject)
		if !isObject {
			continue
		}
		if key == "" {
			sanitized.Set(key, s.sanitizeRootEntry(entry))

			continue
		}
		sanitized.Set(s.maskPath(key), s.sanitizeEntry(entry))
	}

	return sanitized
}

// sanitizeRootEntry keeps less than other entries, because this entry describes the customer's
// own project. Its version and other details aren't read.
func (s *packageLockSanitizer) sanitizeRootEntry(entry *util.OrderedObject) *util.OrderedObject {
	sanitized := util.NewOrderedObject()
	for _, key := range entry.Keys {
		list, isObject := entry.Values[key].(*util.OrderedObject)
		if isObject && util.IsPackageJSONDependencyList(key) {
			sanitized.Set(key, s.sanitizeList(list))
		}
	}

	return sanitized
}

func (s *packageLockSanitizer) sanitizeEntry(entry *util.OrderedObject) *util.OrderedObject {
	sanitized := util.NewOrderedObject()
	isLink, _ := entry.Values["link"].(bool)

	for _, key := range entry.Keys {
		s.setField(sanitized, key, entry.Values[key], isLink)
	}

	return sanitized
}

// setField copies one allow-listed field of a package entry, masked. Other fields are dropped.
func (s *packageLockSanitizer) setField(entry *util.OrderedObject, key string, value any, isLink bool) {
	switch {
	case key == "version":
		if version, isString := value.(string); isString {
			entry.Set(key, util.MaskValue(version))
		}
	case key == "dev" || key == "link":
		if flag, isBool := value.(bool); isBool {
			entry.Set(key, flag)
		}
	case key == "resolved":
		if resolved, isString := value.(string); isString {
			s.setResolved(entry, resolved, isLink)
		}
	case key == "requires" || util.IsPackageJSONDependencyList(key):
		s.setList(entry, key, value)
	}
}

func (s *packageLockSanitizer) setResolved(entry *util.OrderedObject, resolved string, isLink bool) {
	switch {
	case isLink:
		entry.Set("resolved", s.maskPath(resolved))
	case util.MatchesUnsupportedProtocol(resolved):
		entry.Set("resolved", util.MaskValue(resolved))
	}
}

// setList also accepts the boolean that `requires` can be in older lock files.
func (s *packageLockSanitizer) setList(entry *util.OrderedObject, key string, value any) {
	switch typed := value.(type) {
	case *util.OrderedObject:
		entry.Set(key, s.sanitizeList(typed))
	case bool:
		if key == "requires" {
			entry.Set(key, typed)
		}
	}
}

// sanitizeList handles ranges by name as well as the nested entries of version 1.
func (s *packageLockSanitizer) sanitizeList(list *util.OrderedObject) *util.OrderedObject {
	sanitized := util.NewOrderedObject()
	for _, name := range list.Keys {
		switch typed := list.Values[name].(type) {
		case string:
			sanitized.Set(name, util.MaskValue(typed))
		case *util.OrderedObject:
			sanitized.Set(name, s.sanitizeEntry(typed))
		}
	}

	return sanitized
}

// maskPath replaces everything before the package's location in node_modules. A package installed
// under a workspace, such as packages/a/node_modules/qs, stays identifiable as qs.
func (s *packageLockSanitizer) maskPath(path string) string {
	if path == "" {
		return path
	}
	if strings.HasPrefix(path, packageLockModulesPrefix) && !strings.Contains(path, "..") {
		return path
	}

	prefix, modules, hasModules := strings.Cut(path, packageLockModulesInfix)
	if !hasModules || strings.Contains(modules, "..") {
		return s.pathPlaceholder(path)
	}

	return s.pathPlaceholder(prefix) + packageLockModulesInfix + modules
}

func (s *packageLockSanitizer) pathPlaceholder(path string) string {
	if placeholder, found := s.paths[path]; found {
		return placeholder
	}
	placeholder := s.context.PathPlaceholder(packageLockPathFormat, path)
	s.paths[path] = placeholder

	return placeholder
}
