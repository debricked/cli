package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakePathFromManifestFile(t *testing.T) {
	path := MakePathFromManifestFile("pkg/resolution/pm/util/file.json", "file.lock")
	assert.Equal(t, "pkg/resolution/pm/util/file.lock", path)

	path = MakePathFromManifestFile("file.json", "file.lock")
	assert.Equal(t, "file.lock", path)
}
