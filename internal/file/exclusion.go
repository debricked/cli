package file

import (
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
)

func Excluded(exclusions []string, path string) bool {
	for _, exclusion := range exclusions {
		ex := filepath.Clean(exclusion)
		matched, _ := doublestar.PathMatch(ex, path)
		if matched {
			return true
		}
	}

	return false
}
