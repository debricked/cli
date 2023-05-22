package finder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestFindLongestDirMatch(t *testing.T) {

	dirs := []string{"test/", "test2/"}
	file := "test/asd"
	dirMatch, err := findLongestDirMatch(file, dirs)

	assert.Equal(t, "test/", dirMatch)
	assert.Nil(t, err)
}

func TestFindLongestDirMatchErr(t *testing.T) {

	dirs := []string{"test/", "test2/"}
	file := "gest/asd"
	dirMatch, err := findLongestDirMatch(file, dirs)

	assert.Equal(t, "", dirMatch)
	assert.NotNil(t, err)
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
