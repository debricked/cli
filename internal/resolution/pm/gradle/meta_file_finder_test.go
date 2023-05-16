package gradle

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	finder := MetaFileFinder{filepath: FilePath{}}
	paths := []string{filepath.Join("testdata", "project")}
	sMap, gMap, _ := finder.Find(paths)

	assert.Len(t, sMap, 1)
	assert.Len(t, gMap, 1)
}

func TestFindNoFiles(t *testing.T) {
	finder := MetaFileFinder{filepath: FilePath{}}
	paths := []string{filepath.Join("testdata", "project", "subproject")}
	sMap, gMap, _ := finder.Find(paths)

	assert.Len(t, sMap, 0)
	assert.Len(t, gMap, 0)
}

type mockGradleFilePath struct{}

func (m mockGradleFilePath) Walk(root string, walkFn filepath.WalkFunc) error {
	return errors.New("test")
}

func (m mockGradleFilePath) Base(path string) string {
	return filepath.Base(path)
}

func (m mockGradleFilePath) Abs(path string) (string, error) {
	return filepath.Abs(path)
}

func (m mockGradleFilePath) Dir(path string) string {
	return filepath.Dir(path)
}

func TestWalkError(t *testing.T) {
	finder := MetaFileFinder{filepath: mockGradleFilePath{}}
	paths := []string{filepath.Join("testdata", "project", "subproject")}
	_, _, err := finder.Find(paths)
	assert.EqualError(t, err, SetupWalkError{message: "test"}.Error())
}

func TestWalkFuncError(t *testing.T) {
	finder := MetaFileFinder{filepath: FilePath{}}
	paths := []string{filepath.Join("testdata", "test")}
	_, _, err := finder.Find(paths)

	// assert err not nil
	assert.NotNil(t, err)
}
