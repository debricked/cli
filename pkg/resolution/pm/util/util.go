package util

import (
	"fmt"
	"path/filepath"
	"strings"
)

func MakePathFromManifestFile(siblingFile string, fileName string) string {
	dir := filepath.Dir(siblingFile)
	if strings.EqualFold("/", dir) {
		return fmt.Sprintf("/%s", fileName)
	}

	return fmt.Sprintf("%s/%s", dir, fileName)
}
