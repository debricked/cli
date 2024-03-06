package golanfinder

import (
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

func TestFindFiles(t *testing.T) {
	f := GolangFinder{}
	files, err := f.FindFiles([]string{"testdata"}, nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, files)
	assert.Contains(t, files, "testdata/app.go")
	assert.Contains(t, files, "testdata/util.go")

}
