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
		filepath.Join("foo/bar/tree/target/classfiles") + string(os.PathSeparator),
		filepath.Join("foo/bar/tree/asd/hej/target/classfiles") + string(os.PathSeparator),
		filepath.Join("foo/bar/tree/asd/fast/target/classfiles") + string(os.PathSeparator),
	}
	files := []string{
		filepath.Join("foo/bar/tree/pom.xml"),
		filepath.Join("foo/bar/tree/asd/fast/pom.xml"),
	}
	mapFiles := MapFilesToDir(files, dirs)
	fmt.Println(dirs)
	fmt.Println(files)
	fmt.Println(mapFiles)

	assert.Len(t, mapFiles[filepath.Join("foo/bar/tree/pom.xml")], 2)
	assert.Len(t, mapFiles[filepath.Join("foo/bar/tree/asd/fast/pom.xml")], 1)
}

func TestFindPomFileMatch(t *testing.T) {

	classDir := filepath.Join("foo/bar/tree/asd/hej/target/classfiles") + string(os.PathSeparator)
	pomFiles := []string{
		filepath.Join("foo/bar/tree/pom.xml"),
		filepath.Join("foo/bar/tree/asd/fast/pom.xml"),
	}
	pomFileMatch, err := findPomFileMatch(classDir, pomFiles)

	assert.Equal(t, filepath.Join("foo/bar/tree/pom.xml"), pomFileMatch)
	assert.Nil(t, err)
}

func TestFindPomFileMatchDifficult(t *testing.T) {

	classDir := filepath.Join("foo/bar/tree/asd/hej/target/classfiles") + string(os.PathSeparator)
	pomFiles := []string{
		filepath.Join("foo/bar/tree/pom.xml"),
		filepath.Join("foo/bar/tree/asd/fast/pom.xml"),
		filepath.Join("foo/bar/tree/asd/hej/pom.xml"),
	}
	pomFileMatch, err := findPomFileMatch(classDir, pomFiles)

	assert.Equal(t, filepath.Join("foo/bar/tree/asd/hej/pom.xml"), pomFileMatch)
	assert.Nil(t, err)
}

func TestFindPomFileMatchErr(t *testing.T) {

	dirs := []string{"test" + string(os.PathSeparator), "test2" + string(os.PathSeparator)}
	file := filepath.Join("gest", "asd")
	dirMatch, err := findPomFileMatch(file, dirs)

	assert.Equal(t, "", dirMatch)
	assert.NotNil(t, err)
}

func TestMapFilesToDirNoFiles(t *testing.T) {

	dirs := []string{
		filepath.Join("test/asd"),
		filepath.Join("test2/basd/qwe"),
		filepath.Join("test2/test/asd"),
		filepath.Join("test3/tes"),
	}
	files := []string{}
	mapFiles := MapFilesToDir(files, dirs)

	assert.Empty(t, mapFiles)
}

func TestMapFilesToEmptyDir(t *testing.T) {

	dirs := []string{}
	files := []string{
		filepath.Join("test/asd"),
		filepath.Join("test2/basd/qwe"),
		filepath.Join("test2/test/asd"),
		filepath.Join("test3/tes"),
	}
	mapFiles := MapFilesToDir(files, dirs)

	assert.Empty(t, mapFiles)
}

func TestMapFilesToDirNoMatches(t *testing.T) {

	dirs := []string{"tset/pom.xml"}
	files := []string{
		filepath.Join("test/target/classes"),
	}
	mapFiles := MapFilesToDir(files, dirs)

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
