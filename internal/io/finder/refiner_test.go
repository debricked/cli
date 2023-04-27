package finder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindFiles(t *testing.T) {
	files, err := FindFiles([]string{"."}, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, files)
}

func TestFindFilesErr(t *testing.T) {
	files, err := FindFiles([]string{"totaly-not-a-valid-path-123123123"}, nil)
	assert.NotNil(t, err)
	assert.Empty(t, files)
}

func TestConvertPathsToAbsPaths(t *testing.T) {
	before := "321"
	files, err := ConvertPathsToAbsPaths([]string{before})
	after := files[0]
	assert.Nil(t, err)
	assert.True(t, len(before) < len(after))
}

func TestMapFilesToDir(t *testing.T) {

	dirs := []string{"test/", "test2/"}
	files := []string{"test/asd", "test2/basd/qwe", "test2/test/asd", "test3/tes"}
	mapFiles := MapFilesToDir(dirs, files)

	assert.Len(t, mapFiles["test/"], 1)
	assert.Len(t, mapFiles["test2/"], 2)
}

func TestMapFilesToEmptyDir(t *testing.T) {

	dirs := []string{}
	files := []string{"test/asd", "test2/basd/qwe", "test2/test/asd", "test3/tes"}
	mapFiles := MapFilesToDir(dirs, files)

	assert.Empty(t, mapFiles)
}

func TestGCDPath(t *testing.T) {
	files := []string{"test2/bas", "test2/basd/qwe", "test2/test/asd"}
	res := GCDPath(files)

	assert.Equal(t, res, "test2/")
}
