package bazel

import "github.com/debricked/cli/internal/sanitize/pm/util"

// sanitizeInstallJSON keeps each entry's coord and its edges: directDependencies, or
// dependencies if there are no direct ones, because that is what the backend reads. URLs, mirror
// URLs, file paths and hashes are removed: the URLs can hold private hosts and credentials.
//
// The backend needs dependency_tree, so the newer format without it fails here rather than being
// uploaded and giving no packages.
func sanitizeInstallJSON(content []byte, _ *util.FileContext) ([]byte, error) {
	install, err := util.ParseJSONObject(content)
	if err != nil {
		return nil, util.ErrInvalidFormat
	}
	tree, isObject := install.Values["dependency_tree"].(*util.OrderedObject)
	if !isObject {
		return nil, util.ErrInvalidFormat
	}
	entries, isList := tree.Values["dependencies"].([]any)
	if !isList {
		return nil, util.ErrInvalidFormat
	}

	sanitizedEntries := make([]any, 0, len(entries))
	for _, item := range entries {
		entry, isObject := item.(*util.OrderedObject)
		if !isObject {
			return nil, util.ErrInvalidFormat
		}
		sanitizedEntry, err := sanitizeEntry(entry)
		if err != nil {
			return nil, err
		}
		sanitizedEntries = append(sanitizedEntries, sanitizedEntry)
	}

	sanitizedTree := util.NewOrderedObject()
	sanitizedTree.Set("dependencies", sanitizedEntries)
	sanitized := util.NewOrderedObject()
	sanitized.Set("dependency_tree", sanitizedTree)

	return util.EncodeJSON(sanitized)
}

func sanitizeEntry(entry *util.OrderedObject) (*util.OrderedObject, error) {
	coord, isString := entry.Values["coord"].(string)
	if !isString || coord == "" {
		return nil, util.ErrInvalidFormat
	}
	sanitized := util.NewOrderedObject()
	sanitized.Set("coord", coord)

	for _, key := range []string{"directDependencies", "dependencies"} {
		value, present := entry.Values[key]
		if !present {
			continue
		}
		coords, err := coordList(value)
		if err != nil {
			return nil, err
		}
		sanitized.Set(key, coords)

		break
	}

	return sanitized, nil
}

func coordList(value any) ([]any, error) {
	list, isList := value.([]any)
	if !isList {
		return nil, util.ErrInvalidFormat
	}
	for _, item := range list {
		if _, isString := item.(string); !isString {
			return nil, util.ErrInvalidFormat
		}
	}

	return list, nil
}
