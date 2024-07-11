package file

import (
	"fmt"
	"path/filepath"
)

const (
	StrictAll          = 0
	StrictLockAndPairs = 1
	StrictPairs        = 2
)

type Groups struct {
	groups []*Group
}

// Match `format` with `path` and append result to Groups
func (gs *Groups) Match(format *CompiledFormat, path string, lockfileOnly bool) bool {
	dir, file := filepath.Split(path)

	// If it is not a match, return
	manifestFileMatch := !lockfileOnly && format.MatchFile(file)
	lockFileMatch := format.MatchLockFile(file)

	if !manifestFileMatch && !lockFileMatch {
		return false
	}

	if gs.groupExists(format, manifestFileMatch, lockFileMatch, dir, file) {
		return true
	}

	// Create new Group
	var newG *Group
	if manifestFileMatch {
		newG = NewGroup(path, format, []string{})
	} else {
		newG = NewGroup("", format, []string{path})
	}
	gs.Add(*newG)

	return true
}

func (gs *Groups) groupExists(format *CompiledFormat, matchOnManifestFile bool, matchOnLockFile bool, dir string, file string) bool {
	for _, g := range gs.groups {
		if format != g.CompiledFormat {
			continue
		}
		if matchOnLockFile && g.matchLockFile(file, dir) {
			g.LockFiles = append(g.LockFiles, dir+file)

			return true

		} else if matchOnManifestFile && g.matchManifestFile(file, dir) {
			g.ManifestFile = dir + file

			return true
		}
	}

	return false
}

func (gs *Groups) FilterGroupsByStrictness(strictness int) {
	var groups []*Group

	if strictness == StrictAll {
		return
	}

	for _, group := range gs.groups {
		if !group.HasLockFiles() {
			continue
		}

		if strictness == StrictLockAndPairs || (strictness == StrictPairs && group.HasFile()) {
			groups = append(groups, group)
		}
	}

	if len(groups) == 0 && len(gs.groups) > 0 {
		fmt.Println("The following files and directories were filtered out by strictness flag, resulting in no file matches.")
		for _, group := range gs.groups {
			fmt.Println(group.GetAllFiles())
		}
	}

	gs.groups = groups
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
