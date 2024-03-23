package golangfinder

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGolangFindEntrypoint(t *testing.T) {
	f := GolangFinder{}
	files, err := f.FindRoots([]string{filepath.Join("testdata", "app.go"), filepath.Join("testdata", "util.go")})
	assert.Nil(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, filepath.Join("testdata", "app.go"), files[0])
}

func TestGolangFindEntrypointNoMain(t *testing.T) {
	f := GolangFinder{}
	files, err := f.FindRoots([]string{filepath.Join("testdata", "extrapackage", "extra.go")})
	assert.Nil(t, err)
	assert.Empty(t, files)
}

func TestFindFiles(t *testing.T) {
	f := GolangFinder{}
	files, err := f.FindFiles([]string{"testdata"}, nil, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, files)
	assert.Contains(t, files, filepath.Join("testdata", "app.go"))
	assert.Contains(t, files, filepath.Join("testdata", "util.go"))
}

func TestFindDependencyDirs(t *testing.T) {
	f := GolangFinder{}
	files, err := f.FindDependencyDirs([]string{filepath.Join("testdata", "app.go"), filepath.Join("testdata", "util.go")}, false)
	assert.Nil(t, err)
	assert.Empty(t, files)
}

func TestFindFilesExclusions(t *testing.T) {
	finder := GolangFinder{}
	files, err := finder.FindFiles([]string{"testdata"}, []string{"testdata"}, nil)
	assert.Nil(t, err)
	assert.Empty(t, files)
}

func TestIsMainFileError(t *testing.T) {
	finder := GolangFinder{}
	_, err := finder.isMainFile("nonexistent.go")
	assert.NotNil(t, err)
}

func TestFindRootsFileError(t *testing.T) {
	finder := GolangFinder{}
	_, err := finder.FindRoots([]string{"nonexistent.go"})
	assert.NotNil(t, err)
}
