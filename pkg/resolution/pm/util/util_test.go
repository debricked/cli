package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakePathFromManifestFile(t *testing.T) {
	path := MakePathFromManifestFile("pkg/resolution/pm/util/file.json", "file.lock")
	assert.Equal(t, "pkg/resolution/pm/util/file.lock", path)

	path = MakePathFromManifestFile("file.json", "file.lock")
	assert.Equal(t, "./file.lock", path)

	path = MakePathFromManifestFile("/", "file.lock")
	assert.Equal(t, "/file.lock", path)
}
