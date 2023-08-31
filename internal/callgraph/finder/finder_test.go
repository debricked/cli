package finder

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindMavenRoots(t *testing.T) {

	files := []string{"test/asd/pom.xml", "test2/pom.xml", "test2/test/asd/pom.xml", "test3/tes"}
	f := Finder{}
	roots, err := f.FindMavenRoots(files)

	assert.Nil(t, err)
	assert.Len(t, roots, 0)
}

func TestFindJavaClassDirs(t *testing.T) {
	files := []string{"test/asd/pom.xml", "test2/basd/qwe/asd.class", "test2/test/asd", "test3/tes.jar"}
	f := Finder{}
	files, err := f.FindJavaClassDirs(files, false)

	assert.Nil(t, err)
	assert.Len(t, files, 1)
	gt := filepath.Join("test2", "basd", "qwe")
	assert.Equal(t, files[0], gt)

	files = []string{"test/asd/pom.xml", "test2/basd/qwe/asd.class", "test2/test/asd", "test3/tes.jar"}
	files, err = f.FindJavaClassDirs(files, true)

	assert.Nil(t, err)
	assert.Len(t, files, 2)
}

func TestFindFiles(t *testing.T) {
	f := Finder{}
	files, err := f.FindFiles([]string{"."}, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, files)
}

func TestFindFilesErr(t *testing.T) {
	f := Finder{}
	files, err := f.FindFiles([]string{"totaly-not-a-valid-path-123123123"}, nil)
	assert.NotNil(t, err)
	assert.Empty(t, files)
}
