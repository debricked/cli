package composer

import "github.com/debricked/cli/internal/sanitize/pm/util"

// sanitizeComposerJSON keeps only require and require-dev keys, preserving version constraints.
func sanitizeComposerJSON(content []byte, _ *util.FileContext) ([]byte, error) {
	manifest, err := util.ParseJSONObject(content)
	if err != nil {
		return nil, util.ErrInvalidFormat
	}

	sanitized := util.NewOrderedObject()
	for _, listName := range []string{"require", "require-dev"} {
		requirements, isObject := manifest.Values[listName].(*util.OrderedObject)
		if !isObject {
			continue
		}

		sanitizedList := util.NewOrderedObject()
		for _, dependency := range requirements.Keys {
			sanitizedList.Set(dependency, requirements.Values[dependency])
		}
		if len(sanitizedList.Keys) > 0 {
			sanitized.Set(listName, sanitizedList)
		}
	}

	return util.EncodeJSON(sanitized)
}
