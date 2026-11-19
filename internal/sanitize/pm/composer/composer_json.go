package composer

import "github.com/debricked/cli/internal/sanitize/pm/util"

// composerManifestLists are the lists the backend reads from composer.json: their names decide
// which packages are direct dependencies.
var composerManifestLists = []string{"require", "require-dev"}

// sanitizeComposerJSON keeps only require and require-dev, with their names and version
// constraints as they are. Everything else is dropped: repositories, scripts and config often
// hold internal hosts and credentials.
func sanitizeComposerJSON(content []byte, _ *util.FileContext) ([]byte, error) {
	manifest, err := util.ParseJSONObject(content)
	if err != nil {
		return nil, util.ErrInvalidFormat
	}

	sanitized := util.NewOrderedObject()
	for _, listName := range composerManifestLists {
		value, present := manifest.Values[listName]
		if !present {
			continue
		}
		requirements, isObject := value.(*util.OrderedObject)
		if !isObject {
			return nil, util.ErrInvalidFormat
		}
		constraints, err := keepConstraints(requirements)
		if err != nil {
			return nil, err
		}
		if len(constraints.Keys) > 0 {
			sanitized.Set(listName, constraints)
		}
	}

	return util.EncodeJSON(sanitized)
}

// keepConstraints copies a require list. Every constraint must be a string, so nothing nested can
// be copied.
func keepConstraints(requirements *util.OrderedObject) (*util.OrderedObject, error) {
	sanitized := util.NewOrderedObject()
	for _, name := range requirements.Keys {
		constraint, isString := requirements.Values[name].(string)
		if !isString {
			return nil, util.ErrInvalidFormat
		}
		sanitized.Set(name, constraint)
	}

	return sanitized, nil
}

// anyVersion keeps the names of a require list in composer.lock and sets every constraint to "*".
// The backend reads only the names; each version comes from the package's own entry.
func anyVersion(requirements *util.OrderedObject) *util.OrderedObject {
	sanitized := util.NewOrderedObject()
	for _, name := range requirements.Keys {
		sanitized.Set(name, "*")
	}

	return sanitized
}
