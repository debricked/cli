package golanfinder

import (
	"os"
	"testing"

	"github.com/debricked/cli/internal/callgraph/finder"
	"github.com/stretchr/testify/assert"
)

func TestGolangFinderImplementsFinder(t *testing.T) {
	assert.Implements(t, (*finder.IFinder)(nil), new(GolangFinder))
}

func TestGolangFindEntrypoint(t *testing.T) {
	f := GolangFinder{}
	files, err := f.FindRoots([]string{"testdata/app.go", "testdata/util.go"})
	assert.Nil(t, err)
	assert.Len(t, files, 1)
	assert.Equal(t, "testdata/app.go", files[0])
}

func TestGolangFindEntrypointNoMain(t *testing.T) {
	f := GolangFinder{}
	files, err := f.FindRoots([]string{"testdata/extrapackage/extra.go"})
	assert.Nil(t, err)
	assert.Empty(t, files)
}

func TestFindFiles(t *testing.T) {
	f := GolangFinder{}
	files, err := f.FindFiles([]string{"testdata"}, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, files)
	assert.Contains(t, files, "testdata/app.go")
	assert.Contains(t, files, "testdata/util.go")

}

func TestFindDependencyDirs(t *testing.T) {
	f := GolangFinder{}
	files, err := f.FindDependencyDirs([]string{"testdata/app.go", "testdata/util.go"}, false)
	assert.Nil(t, err)
	assert.Empty(t, files)
}

func TestFindFilesWithErrors(t *testing.T) {
	finder := GolangFinder{}
	_, err := finder.FindFiles([]string{"nonexistent"}, nil)
	assert.Error(t, err)

	tempDir, err := os.MkdirTemp("", "testdir")
	assert.Nil(t, err)
	defer os.RemoveAll(tempDir)   // clean up
	err = os.Chmod(tempDir, 0222) // remove read permissions
	assert.Nil(t, err)
	_, err = finder.FindFiles([]string{tempDir}, nil)
	assert.Error(t, err)

}

func TestFindFilesExclusions(t *testing.T) {
	finder := GolangFinder{}
	files, err := finder.FindFiles([]string{"testdata"}, []string{"testdata"})
	assert.Nil(t, err)
	assert.Empty(t, files)
}
