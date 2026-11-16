package python

import (
	"github.com/debricked/cli/internal/sanitize/pm/util"
	"github.com/pelletier/go-toml/v2"
)

func decodeTOMLLock(content []byte) (map[string]any, error) {
	var lock map[string]any
	if err := toml.Unmarshal(content, &lock); err != nil {
		return nil, util.ErrInvalidFormat
	}

	return lock, nil
}

func encodeTOMLLock(packages []any) ([]byte, error) {
	content, err := toml.Marshal(map[string]any{"package": packages})
	if err != nil {
		return nil, util.ErrInvalidFormat
	}

	return content, nil
}

func packageEntries(lock map[string]any) ([]any, error) {
	packages, ok := lock["package"].([]any)
	if !ok {
		return nil, util.ErrInvalidFormat
	}

	return packages, nil
}

func packageTable(value any) (map[string]any, error) {
	entry, ok := value.(map[string]any)
	if !ok {
		return nil, util.ErrInvalidFormat
	}

	return entry, nil
}
