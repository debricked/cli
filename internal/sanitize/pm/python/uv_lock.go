package python

import "github.com/debricked/cli/internal/sanitize/pm/util"

func sanitizeUVLock(content []byte, context *util.FileContext) ([]byte, error) {
	lock, err := decodeTOMLLock(content)
	if err != nil {
		return nil, err
	}
	packages, err := packageEntries(lock)
	if err != nil {
		return nil, err
	}
	entries, projectNames, err := prepareUVPackages(packages, context)
	if err != nil {
		return nil, err
	}
	sanitized, err := sanitizeUVPackages(entries, projectNames)
	if err != nil {
		return nil, err
	}

	return encodeTOMLLock(sanitized)
}

func prepareUVPackages(packages []any, context *util.FileContext) ([]map[string]any, map[string]string, error) {
	projectNames := map[string]string{}
	entries := make([]map[string]any, len(packages))
	for index, value := range packages {
		entry, err := packageTable(value)
		if err != nil {
			return nil, nil, err
		}
		name, ok := entry["name"].(string)
		if !ok {
			return nil, nil, util.ErrInvalidFormat
		}
		if _, ok = entry["version"]; !ok {
			return nil, nil, util.ErrInvalidFormat
		}
		entries[index] = entry
		if isProjectSource(entry["source"]) {
			projectNames[name] = context.ProjectPlaceholder(name)
		}
	}

	return entries, projectNames, nil
}

func sanitizeUVPackages(entries []map[string]any, projectNames map[string]string) ([]any, error) {
	sanitized := make([]any, 0, len(entries))
	for _, entry := range entries {
		packageEntry, err := sanitizeUVPackage(entry, projectNames)
		if err != nil {
			return nil, err
		}
		sanitized = append(sanitized, packageEntry)
	}

	return sanitized, nil
}

func sanitizeUVPackage(entry map[string]any, projectNames map[string]string) (map[string]any, error) {
	name := entry["name"].(string)
	packageEntry := map[string]any{"name": name, "version": entry["version"]}
	if placeholder, ok := projectNames[name]; ok {
		packageEntry["name"] = placeholder
		packageEntry["version"] = "0.0.0"
	}
	if dependencies, ok := entry["dependencies"]; ok {
		redactedDependencies, err := sanitizeUVDependencies(dependencies, projectNames)
		if err != nil {
			return nil, err
		}
		packageEntry["dependencies"] = redactedDependencies
	}

	return packageEntry, nil
}

func sanitizeUVDependencies(value any, projectNames map[string]string) ([]any, error) {
	dependencies, ok := value.([]any)
	if !ok {
		return nil, util.ErrInvalidFormat
	}
	sanitized := make([]any, 0, len(dependencies))
	for _, value := range dependencies {
		dependency, ok := value.(map[string]any)
		if !ok {
			return nil, util.ErrInvalidFormat
		}
		name, ok := dependency["name"].(string)
		if !ok {
			return nil, util.ErrInvalidFormat
		}
		if placeholder, ok := projectNames[name]; ok {
			name = placeholder
		}
		sanitized = append(sanitized, map[string]any{"name": name})
	}

	return sanitized, nil
}

func isProjectSource(source any) bool {
	sourceTable, ok := source.(map[string]any)
	if !ok {
		return false
	}
	_, virtual := sourceTable["virtual"]
	_, editable := sourceTable["editable"]

	return virtual || editable
}
