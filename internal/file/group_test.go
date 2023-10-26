package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExamplePrint() {
	format := Format{
		ManifestFileRegex: "ManifestFileRegex",
		DocumentationUrl:  "https://debricked.com/docs",
		LockFileRegexes:   []string{},
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
		ManifestFile:   "",
		CompiledFormat: nil,
		LockFiles:      []string{"yarn.lock"},
	}
	assert.False(t, g.HasFile(), "failed to assert that group had no dependency file")

	g = Group{
		ManifestFile:   "package.json",
		CompiledFormat: nil,
		LockFiles:      []string{"yarn.lock"},
	}
	assert.True(t, g.HasFile(), "failed to assert that group had a dependency file")
}

func TestGetAllFiles(t *testing.T) {
	g := NewGroup("package.json", nil, []string{"yarn.lock"})
	assert.Len(t, g.GetAllFiles(), 2, "failed to assert number of files")

	g.LockFiles = []string{}
	assert.Len(t, g.GetAllFiles(), 1, "failed to assert number of files")

	g.ManifestFile = ""
	assert.Len(t, g.GetAllFiles(), 0, "failed to assert number of files")
}

func TestGroupMatchLockFile(t *testing.T) {
	var match bool
	g1 := NewGroup("requirements.txt", nil, []string{})
	match = g1.matchLockFile("requirements.txt.pip.debricked.lock", "")
	assert.Equal(t, match, true)
	g2 := NewGroup("/home/requirements-test.txt", nil, []string{})
	match = g2.matchLockFile("requirements-test.txt.pip.debricked.lock", "/home/")
	assert.Equal(t, match, true)
	g3 := NewGroup("requirements.txt", nil, []string{})
	match = g3.matchLockFile("requirements.dev.txt.pip.debricked.lock", "")
	assert.Equal(t, match, false)
	g4 := NewGroup("requirements-test.txt", nil, []string{})
	match = g4.matchLockFile("requirements.txt.pip.debricked.lock", "")
	assert.Equal(t, match, false)
	g5 := NewGroup("requirements-test.txt", nil, []string{})
	match = g5.matchLockFile("requirements.txt-test.txt.pip.debricked.lock", "")
	assert.Equal(t, match, false)

	// Check that match fails if different directories
	g6 := NewGroup("requirements.txt", nil, []string{})
	match = g6.matchLockFile("requirements.txt.pip.debricked.lock", "some/other/directory")
	assert.Equal(t, match, false)

	// Check that match fails if there is no manifest file - functionality may change in the future.
	g7 := NewGroup("", nil, []string{})
	match = g7.matchLockFile("requirements.txt.pip.debricked.lock", "some/other/directory")
	assert.Equal(t, match, false)
}

func TestGroupMatchManifestFile(t *testing.T) {
	var match bool
	g1 := NewGroup("", nil, []string{"/home/requirements.txt.pip.debricked.lock"})
	match = g1.matchManifestFile("requirements.txt", "/home/")
	assert.Equal(t, match, true)

	// Check that match fails if different directories
	g6 := NewGroup("", nil, []string{"requirements.txt.pip.debricked.lock"})
	match = g6.matchManifestFile("requirements.txt", "some/other/directory")
	assert.Equal(t, match, false)

	// Check that match fails when manifest file exists -- only one per group for now.
	g7 := NewGroup("requirements.txt", nil, []string{})
	match = g7.matchManifestFile("requirements.dev.txt", "")
	assert.Equal(t, match, false)
}

func TestGroupMatchFile(t *testing.T) {
	var match bool
	match = matchFile("package.json", "package-lock.json")
	assert.Equal(t, match, true)
	match = matchFile("requirements.txt", "requirements.txt.pip.debricked.lock")
	assert.Equal(t, match, true)
	match = matchFile("requirements.txt", "requirements.txt.pip.debricked.lock")
	assert.Equal(t, match, true)
	match = matchFile("requirements-test.txt", "requirements.txt.pip.debricked.lock")
	assert.Equal(t, match, false)
}
