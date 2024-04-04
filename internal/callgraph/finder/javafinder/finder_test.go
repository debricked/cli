package javafinder

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindMavenRoots(t *testing.T) {

	files := []string{"test/asd/pom.xml", "test2/pom.xml", "test2/test/asd/pom.xml", "test3/tes"}
	f := JavaFinder{}
	roots, err := f.FindRoots(files)

	assert.Nil(t, err)
	assert.Len(t, roots, 0)
}

func TestFindDependencyDirs(t *testing.T) {
	files := []string{"test/asd/pom.xml", "test2/basd/qwe/asd.class", "test2/test/asd", "test3/tes.jar"}
	f := JavaFinder{}
	files, err := f.FindDependencyDirs(files, false)

	assert.Nil(t, err)
	assert.Len(t, files, 1)
	gt := filepath.Join("test2", "basd", "qwe")
	assert.Equal(t, files[0], gt)

	files = []string{"test/asd/pom.xml", "test2/basd/qwe/asd.class", "test2/test/asd", "test3/tes.jar"}
	files, err = f.FindDependencyDirs(files, true)

	assert.Nil(t, err)
	assert.Len(t, files, 2)
}

func TestFindFiles(t *testing.T) {
	f := JavaFinder{}
	files, err := f.FindFiles([]string{"."}, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, files)
}

func TestFindFilesErr(t *testing.T) {
	f := JavaFinder{}
	files, err := f.FindFiles([]string{"totaly-not-a-valid-path-123123123"}, nil)
	assert.NotNil(t, err)
	assert.Empty(t, files)
}

func TestFindFilesExcluded(t *testing.T) {
	f := JavaFinder{}
	project_path, err := filepath.Abs("testdata/test_project")
	assert.Nil(t, err)
	files, err := f.FindFiles([]string{project_path}, nil)
	assert.Nil(t, err)
	assert.Len(t, files, 2)
	files, err = f.FindFiles([]string{project_path}, []string{"**/excluded_folder/**"})
	assert.Nil(t, err)
	assert.Len(t, files, 1)
	files, err = f.FindFiles([]string{project_path}, []string{"excluded_folder"})
	assert.Nil(t, err)
	assert.Len(t, files, 2)
	files, err = f.FindFiles([]string{project_path}, []string{"**/excluded*/**"})
	assert.Nil(t, err)
	assert.Len(t, files, 1)
	files, err = f.FindFiles([]string{project_path}, []string{"**/excluded_file.txt"})
	assert.Nil(t, err)
	assert.Len(t, files, 1)
}
