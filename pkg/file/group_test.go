package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func ExamplePrint() {
	format := Format{
		Regex:            "Regex",
		DocumentationUrl: "https://debricked.com/docs",
		LockFileRegexes:  []string{},
	}
	compiledFormat, _ := NewCompiledFormat(&format)
	g := NewGroup("package.json", compiledFormat, []string{"yarn.lock"})
	g.Print()
	// output:
	// package.json
	//  * yarn.lock
}

func TestHasFile(t *testing.T) {
	g := Group{
		FilePath:       "",
		CompiledFormat: nil,
		RelatedFiles:   []string{"yarn.lock"},
	}
	assert.False(t, g.HasFile(), "failed to assert that group had no dependency file")

	g = Group{
		FilePath:       "package.json",
		CompiledFormat: nil,
		RelatedFiles:   []string{"yarn.lock"},
	}
	assert.True(t, g.HasFile(), "failed to assert that group had a dependency file")
}

func TestGetAllFiles(t *testing.T) {
	g := NewGroup("package.json", nil, []string{"yarn.lock"})
	assert.Len(t, g.GetAllFiles(), 2, "failed to assert number of files")

	g.RelatedFiles = []string{}
	assert.Len(t, g.GetAllFiles(), 1, "failed to assert number of files")

	g.FilePath = ""
	assert.Len(t, g.GetAllFiles(), 0, "failed to assert number of files")
}
