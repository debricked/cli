package pub

import (
	"github.com/debricked/cli/internal/sanitize/pm/util"
	"gopkg.in/yaml.v3"
)

// pubspecLockFields are the fields of a package that the backend reads. description and source
// are removed: they hold git URLs, local paths and private hosts.
var pubspecLockFields = []string{"dependency", "version"}

// sanitizePubspecLock keeps the name, dependency type and version of each package.
func sanitizePubspecLock(content []byte, _ *util.FileContext) ([]byte, error) {
	root, err := util.ParseYAMLMapping(content)
	if err != nil {
		return nil, err
	}
	packages := util.YAMLValue(root, "packages")
	if packages == nil || packages.Kind != yaml.MappingNode {
		return nil, util.ErrInvalidFormat
	}

	sanitizedPackages := util.NewOrderedObject()
	util.ForEachYAMLPair(packages, func(name string, entry *yaml.Node) {
		if err != nil {
			return
		}
		if entry.Kind != yaml.MappingNode {
			err = util.ErrInvalidFormat

			return
		}
		sanitizedEntry := util.NewOrderedObject()
		for _, field := range pubspecLockFields {
			if value, isString := util.YAMLString(util.YAMLValue(entry, field)); isString {
				sanitizedEntry.Set(field, value)
			}
		}
		sanitizedPackages.Set(name, sanitizedEntry)
	})
	if err != nil {
		return nil, err
	}

	sanitized := util.NewOrderedObject()
	sanitized.Set("packages", sanitizedPackages)

	return util.EncodeYAML(sanitized), nil
}
