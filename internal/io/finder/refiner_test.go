package finder

import (
	"fmt"
	"os"
	"path/filepath"
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

	dirs := []string{
		"test" + string(os.PathSeparator),
		"test2" + string(os.PathSeparator),
	}
	files := []string{
		filepath.Join("test/asd"),
		filepath.Join("test2/basd/qwe"),
		filepath.Join("test2/test/asd"),
		filepath.Join("test3/tes"),
	}
	mapFiles := MapFilesToDir(dirs, files)
	fmt.Println(dirs)
	fmt.Println(files)

	assert.Len(t, mapFiles["test"+string(os.PathSeparator)], 1)
	assert.Len(t, mapFiles["test2"+string(os.PathSeparator)], 2)
}

func TestFindLongestDirMatch(t *testing.T) {

	dirs := []string{"test" + string(os.PathSeparator), "test2" + string(os.PathSeparator)}
	file := filepath.Join("test", "asd")
	dirMatch, err := findLongestDirMatch(file, dirs)

	assert.Equal(t, "test"+string(os.PathSeparator), dirMatch)
	assert.Nil(t, err)
}

func TestFindLongestDirMatchErr(t *testing.T) {

	dirs := []string{"test" + string(os.PathSeparator), "test2" + string(os.PathSeparator)}
	file := filepath.Join("gest", "asd")
	dirMatch, err := findLongestDirMatch(file, dirs)

	assert.Equal(t, "", dirMatch)
	assert.NotNil(t, err)
}

func TestMapFilesToEmptyDir(t *testing.T) {

	dirs := []string{}
	files := []string{
		filepath.Join("test/asd"),
		filepath.Join("test2/basd/qwe"),
		filepath.Join("test2/test/asd"),
		filepath.Join("test3/tes"),
	}
	mapFiles := MapFilesToDir(dirs, files)

	assert.Empty(t, mapFiles)
}

func TestGCDPath(t *testing.T) {
	files := []string{
		filepath.Join("test2/bas"),
		filepath.Join("test2/basd/qwe"),
		filepath.Join("test2/test/asd"),
	}
	res := GCDPath(files)

	assert.Equal(t, res, "test2"+string(os.PathSeparator))
}

func TestNoGCDPath(t *testing.T) {
	files := []string{
		filepath.Join("nogcdtest2/bas"),
		filepath.Join("test2/basd/qwe"),
	}
	res := GCDPath(files)

	assert.Equal(t, res, "")
}
