package pub

import (
	"github.com/debricked/cli/internal/sanitize/pm/util"
	"gopkg.in/yaml.v3"
)

// anyVersion replaces a dependency whose source is a git repository, a path or a private host.
const anyVersion = "any"

// pubspecDependencyLists decide each package's scope. The lock file alone isn't enough: a dev
// dependency that is also overridden is "direct overridden" there, and would count as production.
var pubspecDependencyLists = []string{"dependencies", "dev_dependencies", "dependency_overrides"}

// sanitizePubspecYAML keeps only the dependency lists. A plain version constraint is kept, a
// table keeps only its version, and anything else becomes "any".
func sanitizePubspecYAML(content []byte, _ *util.FileContext) ([]byte, error) {
	root, err := util.ParseYAMLMapping(content)
	if err != nil {
		return nil, err
	}

	sanitized := util.NewOrderedObject()
	for _, listName := range pubspecDependencyLists {
		list := util.YAMLValue(root, listName)
		if list == nil || util.IsYAMLNull(list) {
			continue
		}
		if list.Kind != yaml.MappingNode {
			return nil, util.ErrInvalidFormat
		}
		sanitizedList := util.NewOrderedObject()
		util.ForEachYAMLPair(list, func(name string, value *yaml.Node) {
			sanitizedList.Set(name, versionConstraint(value))
		})
		sanitized.Set(listName, sanitizedList)
	}
	if len(sanitized.Keys) == 0 {
		sanitized.Set("dependencies", util.NewOrderedObject())
	}

	return util.EncodeYAML(sanitized), nil
}

func versionConstraint(value *yaml.Node) string {
	if value.Kind == yaml.MappingNode {
		value = util.YAMLValue(value, "version")
	}
	if version, isString := util.YAMLString(value); isString && util.IsPlainVersion(version) {
		return version
	}

	return anyVersion
}
