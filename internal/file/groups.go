package file

import (
	"path/filepath"
	"strings"
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
	fileMatch := !lockfileOnly && format.MatchFile(file)
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

func matchLockToManifest(manifestGroup, lock string) bool {
	var matchName bool
	if strings.HasPrefix(manifestGroup, "composer") {
		// Naive first implementation assuming projects use composer.json and composer.lock.
		matchName = strings.HasPrefix(lock, "composer")
	} else {
		matchName = strings.HasPrefix(lock, manifestGroup)
	}

	return matchName
}

func matchManifestToLock(lockGroup, file string) bool {
	var matchName bool
	if strings.HasPrefix(lockGroup, "composer") {
		// Naive first implementation assuming projects use composer.json and composer.lock.
		matchName = strings.HasPrefix(file, "composer")
	} else {
		matchName = strings.HasPrefix(lockGroup, file)
	}

	return matchName
}

func (gs *Groups) matchExistingGroup(format *CompiledFormat, fileMatch bool, lockFileMatch bool, dir string, file string) bool {
	for _, g := range gs.groups {
		if format != g.CompiledFormat {
			continue
		}
		var groupDir, manifestFile, lockFile string
		var fileMatch, manifestMatchLock, lockMatchManifest bool
		if g.HasFile() {
			// Group has a manifest file.
			// Check if the input file is a lock file for the manifest file.
			groupDir, manifestFile = filepath.Split(g.ManifestFile)
			lockMatchManifest = matchLockToManifest(manifestFile, file)
		} else {
			// Group does not have a manifest file.
			// Check if the input file is a manifest file for an existing lock file.
			groupDir, lockFile = filepath.Split(g.LockFiles[0])
			manifestMatchLock = matchManifestToLock(lockFile, file)

		}
		fileMatch = manifestMatchLock || lockMatchManifest
		correctFilePath := g.checkFilePathDependantCases(fileMatch, lockFileMatch, file)
		if !fileMatch || groupDir != dir || !correctFilePath {
			continue
		}
		if manifestMatchLock {
			g.ManifestFile = dir + file

			return true
		} else if lockMatchManifest {
			g.LockFiles = append(g.LockFiles, dir+file)

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
