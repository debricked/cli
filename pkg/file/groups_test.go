package file

import (
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
			if gs.Size() != 0 {
				t.Error("failed to assert that groups size was 0")
			}

			gs.Match(compiledF, c.sequence[0], false)
			if gs.Size() != 1 {
				t.Error("failed to assert that there was one Group in Groups")
			}
			gs.Match(compiledF, c.sequence[1], false)
			if gs.Size() != 1 {
				t.Error("failed to assert that there was one Group in Groups")
			}

			g := gs.groups[0]
			if g.FilePath != c.file {
				t.Error("failed to assert that FilePath had correct value directory/composer.json")
			}

			if len(g.RelatedFiles) != 1 {
				t.Error("failed to assert that there was one lock file")
			}

			lockFile := g.RelatedFiles[0]
			if lockFile != c.lockFile {
				t.Error("failed to assert lock file name")
			}
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
