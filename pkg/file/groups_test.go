package file

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatch(t *testing.T) {
	f := Format{
		Regex:            "composer\\.json",
		DocumentationUrl: "",
		LockFileRegexes:  []string{"composer\\.lock"},
	}
	compiledF, _ := NewCompiledFormat(&f)
	var gs Groups

	cases := []struct {
		name     string
		file     string
		lockFile string
		sequence []string
	}{
		{
			name:     "FileFirst",
			file:     "directory/composer.json",
			lockFile: "directory/composer.lock",
			sequence: []string{"directory/composer.json", "directory/composer.lock"},
		},
		{
			name:     "LockFileFirst",
			file:     "directory/composer.json",
			lockFile: "directory/composer.lock",
			sequence: []string{"directory/composer.lock", "directory/composer.json"},
		},
	}

	for _, c := range cases {
		gs.groups = nil
		t.Run(c.name, func(t *testing.T) {
			gs.Match(compiledF, "/", false)
			assert.Equal(t, 0, gs.Size())

			gs.Match(compiledF, c.sequence[0], false)
			assert.Equal(t, 1, gs.Size(), "failed to assert that there was one Group in Groups")

			gs.Match(compiledF, c.sequence[1], false)
			assert.Equal(t, 1, gs.Size(), "failed to assert that there was one Group in Groups")

			g := gs.groups[0]
			assert.Equal(t, c.file, g.FilePath, "failed to assert that FilePath had correct value directory/composer.json")

			assert.Len(t, g.RelatedFiles, 1, "failed to assert that there was one lock file")

			lockFile := g.RelatedFiles[0]
			assert.Equal(t, c.lockFile, lockFile, "failed to assert lock file name")
		})
	}
}

func TestGetFiles(t *testing.T) {
	g1 := NewGroup("file1", nil, []string{"lockfile1"})
	g2 := NewGroup("", nil, []string{"lockfile2"})

	gs := Groups{}
	gs.Add(*g1)
	gs.Add(*g2)
	files := gs.GetFiles()
	const nbrOfFiles = 3
	if len(files) != nbrOfFiles {
		t.Errorf("failed to assert that there was %d files", nbrOfFiles)
	}
}

func TestFilterGroupsByStrictness(t *testing.T) {
	g1 := NewGroup("file1", nil, []string{})
	g2 := NewGroup("", nil, []string{"lockfile2"})
	g3 := NewGroup("file3", nil, []string{"lockfile3"})

	gs := Groups{}
	gs.Add(*g1)
	gs.Add(*g2)
	gs.Add(*g3)

	cases := []struct {
		name                   string
		strictness             int
		expectedNumberOfGroups int
		expectedManifestFile   string
		expectedLockFiles      []string
	}{
		{
			name:                   "StrictnessSetTo0",
			strictness:             StrictAll,
			expectedNumberOfGroups: 3,
			expectedManifestFile:   "file1",
			expectedLockFiles:      []string{},
		},
		{
			name:                   "StrictnessSetTo1",
			strictness:             StrictLockAndPairs,
			expectedNumberOfGroups: 2,
			expectedManifestFile:   "",
			expectedLockFiles: []string{
				"lockfile2",
			},
		},
		{
			name:                   "StrictnessSetTo2",
			strictness:             StrictPairs,
			expectedNumberOfGroups: 1,
			expectedManifestFile:   "file3",
			expectedLockFiles: []string{
				"lockfile3",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			gs.FilterGroupsByStrictness(c.strictness)
			fileGroup := gs.groups[0]

			assert.Equalf(
				t,
				c.expectedNumberOfGroups,
				gs.Size(),
				"failed to assert that %d groups were created. %d were found",
				c.expectedNumberOfGroups,
				gs.Size(),
			)
			assert.Equalf(
				t,
				fileGroup.FilePath,
				c.expectedManifestFile,
				"actual manifest file %s doesn't match expected %s",
				fileGroup.FilePath,
				c.expectedManifestFile,
			)
			assert.EqualValuesf(
				t,
				fileGroup.RelatedFiles,
				c.expectedLockFiles,
				"actual lock files list %s doesn't match expected %s",
				fileGroup.RelatedFiles,
				c.expectedLockFiles,
			)
		})
	}
}
