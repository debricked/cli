package util

import (
	"path/filepath"
)

func MakePathFromManifestFile(siblingFile string, fileName string) string {
	dir := filepath.Dir(siblingFile)

	return filepath.Join(dir, fileName)
}

func MakePathFromManifestFileExtension(siblingFile string, fileName string) string {
	dir := filepath.Dir(siblingFile)
	name := filepath.Base(fileName)

	return filepath.Join(dir, name)
}
