package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func MakePathFromManifestFile(siblingFile string, fileName string) string {
	dir := filepath.Dir(siblingFile)
	if strings.EqualFold(string(os.PathSeparator), dir) {
		return fmt.Sprintf("%s%s", string(os.PathSeparator), fileName)
	}

	return fmt.Sprintf("%s%s%s", dir, string(os.PathSeparator), fileName)
}

func MakePathFromManifestFileExtension(siblingFile string, fileName string) string {
	dir := filepath.Dir(siblingFile)
	name := filepath.Base(fileName)

	return filepath.Join(dir, name)
}
