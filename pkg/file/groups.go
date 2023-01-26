package file

import "path/filepath"

type Groups struct {
	groups []*Group
}

// Match `format` with `path` and append result to Groups
func (gs *Groups) Match(format *CompiledFormat, path string, lockfileOnly bool) bool {
	dir, file := filepath.Split(path)

	// If it is not a match, return
	fileMatch := format.MatchFile(file)
	if lockfileOnly {
		fileMatch = false
	}

	lockFileMatch := format.MatchLockFile(file)

	if !fileMatch && !lockFileMatch {
		return false
	}

	matched := gs.matchExistingGroup(format, fileMatch, lockFileMatch, dir, file)
	if matched {
		return true
	}

	// Create new Group
	var newG *Group
	if fileMatch {
		newG = NewGroup(path, format, []string{})
	} else {
		newG = NewGroup("", format, []string{path})
	}
	gs.Add(*newG)

	return true
}

func (gs *Groups) matchExistingGroup(format *CompiledFormat, fileMatch bool, lockFileMatch bool, dir string, file string) bool {
	for _, g := range gs.groups {
		var gDir string
		if g.HasFile() {
			gDir, _ = filepath.Split(g.FilePath)
		} else {
			gDir, _ = filepath.Split(g.RelatedFiles[0])
		}

		if gDir == dir && format == g.CompiledFormat {
			if fileMatch {
				g.FilePath = dir + file

				return true
			} else if lockFileMatch {
				g.RelatedFiles = append(g.RelatedFiles, dir+file)

				return true
			}
		}
	}

	return false
}

func (gs *Groups) ToSlice() []Group {
	var groups []Group
	for _, g := range gs.groups {
		groups = append(groups, *g)
	}

	return groups
}

func (gs *Groups) Size() int {
	return len(gs.groups)
}

func (gs *Groups) Add(g Group) {
	gs.groups = append(gs.groups, &g)
}

func (gs *Groups) GetFiles() []string {
	var files []string
	for _, g := range gs.groups {
		files = append(files, g.GetAllFiles()...)
	}

	return files
}
