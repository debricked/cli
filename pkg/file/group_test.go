package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestCheckFilePathDependantCasesEmptyGroup(t *testing.T) {
	group := NewGroup("", nil, []string{})
	var check bool

	check = group.checkFilePathDependantCases(true, false, "requirements.txt")
	assert.True(t, check)

	check = group.checkFilePathDependantCases(true, false, "requirements-dev.txt")
	assert.True(t, check)

	check = group.checkFilePathDependantCases(true, false, "requirements-dev.txt")
	assert.True(t, check)

	check = group.checkFilePathDependantCases(false, true, "requirements.txt.pip.debricked.lock")
	assert.False(t, check)
}

func TestCheckFilePathDependantCasesWithFilePath(t *testing.T) {
	group := NewGroup("requirements.txt", nil, []string{})
	var check bool

	check = group.checkFilePathDependantCases(false, true, "requirements.txt.pip.debricked.lock")
	assert.True(t, check)

	check = group.checkFilePathDependantCases(false, true, "requirements-dev.txt.pip.debricked.lock")
	assert.False(t, check)
}

func TestCheckFilePathDependantCasesWithLockFile(t *testing.T) {
	group := NewGroup("", nil, []string{"requirements.txt.pip.debricked.lock"})
	var check bool

	check = group.checkFilePathDependantCases(true, false, "requirements.txt")
	assert.True(t, check)

	check = group.checkFilePathDependantCases(true, false, "requirements-dev.txt")
	assert.False(t, check)

	check = group.checkFilePathDependantCases(true, false, "requirements-dev.txt")
	assert.False(t, check)

	check = group.checkFilePathDependantCases(false, false, "file.txt")
	assert.True(t, check)
}
