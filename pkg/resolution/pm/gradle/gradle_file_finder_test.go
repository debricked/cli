package gradle

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindGradleProjectFiles(t *testing.T) {
	finder := FileFinder{filepath: FilePath{}}
	paths := []string{filepath.Join("testdata", "project")}
	sMap, gMap, _ := finder.FindGradleProjectFiles(paths)

	assert.Len(t, sMap, 1)
	assert.Len(t, gMap, 1)
}

func TestFindGradleProjectFilesNoFiles(t *testing.T) {
	finder := FileFinder{filepath: FilePath{}}
	paths := []string{filepath.Join("testdata", "project", "subproject")}
	sMap, gMap, _ := finder.FindGradleProjectFiles(paths)

	assert.Len(t, sMap, 0)
	assert.Len(t, gMap, 0)
}

type mockFilePath struct{}

func (m mockFilePath) Walk(root string, walkFn filepath.WalkFunc) error {
	return errors.New("test")
}

func (m mockFilePath) Base(path string) string {
	return filepath.Base(path)
}

func (m mockFilePath) Abs(path string) (string, error) {
	return filepath.Abs(path)
}

func (m mockFilePath) Dir(path string) string {
	return filepath.Dir(path)
}

func TestWalkError(t *testing.T) {
	finder := FileFinder{filepath: mockFilePath{}}
	paths := []string{filepath.Join("testdata", "project", "subproject")}
	_, _, err := finder.FindGradleProjectFiles(paths)
	assert.EqualError(t, err, GradleSetupWalkError{message: "test"}.Error())
}

func TestWalkFuncError(t *testing.T) {
	finder := FileFinder{filepath: FilePath{}}
	paths := []string{filepath.Join("testdata", "test")}
	_, _, err := finder.FindGradleProjectFiles(paths)

	// assert err not nil
	assert.NotNil(t, err)
}
