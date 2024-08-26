package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	f := Format{
		ManifestFileRegex: "composer\\.json",
		DocumentationUrl:  "",
		LockFileRegexes:   []string{"composer\\.lock"},
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
			assert.Equal(t, c.file, g.ManifestFile, "failed to assert that ManifestFile had correct value directory/composer.json")

			assert.Len(t, g.LockFiles, 1, "failed to assert that there was one lock file")

			lockFile := g.LockFiles[0]
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
				fileGroup.ManifestFile,
				c.expectedManifestFile,
				"actual manifest file %s doesn't match expected %s",
				fileGroup.ManifestFile,
				c.expectedManifestFile,
			)
			assert.EqualValuesf(
				t,
				fileGroup.LockFiles,
				c.expectedLockFiles,
				"actual lock files list %s doesn't match expected %s",
				fileGroup.LockFiles,
				c.expectedLockFiles,
			)
		})
	}
}

func TestMatchGroupsExpected(t *testing.T) {
	setUp(true)

	testData := map[string][]string{
		"foo/bar/cloud/package.json":                        {"foo/bar/cloud/yarn.lock"},
		"foo/bar/examples/test/requirements.txt":            {},
		"foo/asd/requirements-test-dev.txt":                 {"foo/asd/.requirements-test-dev.txt.pip.debricked.lock"},
		"foo/asd/requirements-test.txt":                     {"foo/asd/.requirements-test.txt.pip.debricked.lock"},
		"foo/asd/requirements.txt":                          {"foo/asd/requirements.txt.pip.debricked.lock"},
		"foo/asd/requirements-api.txt":                      {},
		"foo/asd/src/main/event_listeners/requirements.txt": {"foo/asd/src/main/event_listeners/requirements.txt.pip.debricked.lock"},
		"foo/asd/src/main/util/test/composer.json":          {},
	}

	var groups Groups
	lockfileOnly := false
	formats, _ := finder.GetSupportedFormats()
	for key, values := range testData {
		paths := append(values, key)
		for _, path := range paths {
			for _, format := range formats {
				if groups.Match(format, path, lockfileOnly) {

					break
				}
			}
		}
	}

	for _, group := range groups.groups {
		assert.Equal(t, testData[group.ManifestFile], group.LockFiles)
	}

	assert.Equal(t, len(testData), len(groups.groups))
}

func TestGroupsMatchWorkspaces(t *testing.T) {
	g1 := NewGroup("package.json", nil, []string{"package-lock.json"})
	g2 := NewGroup("", nil, []string{"lockfile2"})
	g3 := NewGroup("package/file3", nil, []string{})
	g4 := NewGroup("pack/file3", nil, []string{})

	gs := Groups{}
	gs.Add(*g1)
	gs.Add(*g2)
	gs.Add(*g3)
	gs.Add(*g4)

	workspaceManifest := WorkspaceManifest{
		LockFile:     "package-lock.json",
		RootManifest: "package.json",
		Workspaces:   []string{"package/*"},
	}
	err := gs.matchWorkspace(workspaceManifest)
	assert.NoError(t, err)
	for _, g := range gs.groups {
		if g.ManifestFile == "package/file3" {
			assert.Equal(t, g.LockFiles[0], g1.LockFiles[0])
		}
		if g.ManifestFile == "pack/file3" {
			assert.Equal(t, len(g3.LockFiles), 0)
		}
	}

}

func TestAddWorkspaceLockFiles(t *testing.T) {
	g1 := NewGroup("testdata/workspace/package.json", nil, []string{"testdata/workspace/package-lock.json"})
	g2 := NewGroup("testdata/workspace/packages/package_one/package.json", nil, []string{})
	g3 := NewGroup("testdata/workspace/packages/package_two/package.json", nil, []string{})

	gs := Groups{}
	gs.Add(*g1)
	gs.Add(*g2)
	gs.Add(*g3)
	err := gs.addWorkspaceLockFiles()
	assert.NoError(t, err)
	for _, g := range gs.groups {
		assert.Equal(t, len(g.LockFiles), 1)
		if len(g.LockFiles) == 1 {
			assert.Equal(t, g.LockFiles[0], g1.LockFiles[0])
		}
	}
}
