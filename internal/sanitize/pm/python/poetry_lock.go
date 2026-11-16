package python

import "github.com/debricked/cli/internal/sanitize/pm/util"

func sanitizePoetryLock(content []byte, _ *util.FileContext) ([]byte, error) {
	lock, err := decodeTOMLLock(content)
	if err != nil {
		return nil, err
	}
	packages, err := packageEntries(lock)
	if err != nil {
		return nil, err
	}

	sanitized := make([]any, 0, len(packages))
	for _, value := range packages {
		entry, err := packageTable(value)
		if err != nil {
			return nil, err
		}
		packageEntry := map[string]any{}
		if name, ok := entry["name"]; ok {
			packageEntry["name"] = name
		} else {
			return nil, util.ErrInvalidFormat
		}
		if version, ok := entry["version"]; ok {
			packageEntry["version"] = version
		} else {
			return nil, util.ErrInvalidFormat
		}
		if dependencies, ok := entry["dependencies"]; ok {
			dependencyTable, ok := dependencies.(map[string]any)
			if !ok {
				return nil, util.ErrInvalidFormat
			}
			redactedDependencies := make(map[string]any, len(dependencyTable))
			for name := range dependencyTable {
				redactedDependencies[name] = "*"
			}
			packageEntry["dependencies"] = redactedDependencies
		}
		sanitized = append(sanitized, packageEntry)
	}

	return encodeTOMLLock(sanitized)
}
