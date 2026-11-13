package npm

import "github.com/debricked/cli/internal/sanitize/pm/util"

// sanitizePackageJSON keeps only the dependency lists, with masked values.
func sanitizePackageJSON(content []byte, _ *util.FileContext) ([]byte, error) {
	manifest, err := util.ParseJSONObject(content)
	if err != nil {
		return nil, util.ErrInvalidFormat
	}

	sanitized := util.NewOrderedObject()
	for _, listName := range util.PackageJSONDependencyLists {
		list, isObject := manifest.Values[listName].(*util.OrderedObject)
		if !isObject {
			continue
		}
		maskedList := util.NewOrderedObject()
		for _, name := range list.Keys {
			if value, isString := list.Values[name].(string); isString {
				maskedList.Set(name, util.MaskValue(value))
			}
		}
		sanitized.Set(listName, maskedList)
	}

	return util.EncodeJSON(sanitized)
}
