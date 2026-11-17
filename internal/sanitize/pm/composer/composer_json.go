package composer

import "github.com/debricked/cli/internal/sanitize/pm/util"

// composerManifestLists are the lists the backend reads from composer.json: their names decide
// which packages are direct dependencies.
var composerManifestLists = []string{"require", "require-dev"}

// sanitizeComposerJSON keeps only the package names in require and require-dev. Everything else
// is dropped: repositories, scripts and config often hold internal hosts and credentials.
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
		if len(requirements.Keys) > 0 {
			sanitized.Set(listName, anyVersion(requirements))
		}
	}

	return util.EncodeJSON(sanitized)
}

// anyVersion keeps the names of a require list and sets every constraint to "*". The backend reads
// only the names; each version comes from the package's own entry in composer.lock. A constraint
// can hold a branch name, such as dev-feature/PROJ-123-new-login, so it is never copied.
func anyVersion(requirements *util.OrderedObject) *util.OrderedObject {
	sanitized := util.NewOrderedObject()
	for _, name := range requirements.Keys {
		sanitized.Set(name, "*")
	}

	return sanitized
}
