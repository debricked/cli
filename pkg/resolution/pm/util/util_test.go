package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestMakePathFromManifestFile(t *testing.T) {
	manifestFile := filepath.Join("pkg", "resolution", "pm", "util", "file.json")
	path := MakePathFromManifestFile(manifestFile, "file.lock")
	lockFile := filepath.Join("pkg", "resolution", "pm", "util", "file.lock")

	assert.Equal(t, lockFile, path)

	path = MakePathFromManifestFile("file.json", "file.lock")
	lockFile = fmt.Sprintf(".%s%s", string(os.PathSeparator), "file.lock")
	assert.Equal(t, lockFile, path)

	path = MakePathFromManifestFile(string(os.PathSeparator), "file.lock")
	assert.Equal(t, fmt.Sprintf("%s%s", string(os.PathSeparator), "file.lock"), path)
}
