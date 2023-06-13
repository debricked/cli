package file

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Group struct {
	ManifestFile   string          `json:"manifestFile"`
	CompiledFormat *CompiledFormat `json:"-"`
	LockFiles      []string        `json:"lockFiles"`
}

func NewGroup(manifestFile string, format *CompiledFormat, lockFiles []string) *Group {
	return &Group{ManifestFile: manifestFile, CompiledFormat: format, LockFiles: lockFiles}
}

func (fileGroup *Group) Print() {
	hasFile := fileGroup.HasFile()
	if hasFile {
		fmt.Println(fileGroup.ManifestFile)
	}
	for _, filePath := range fileGroup.LockFiles {
		if hasFile {
			fmt.Println(" * " + filePath)
		} else {
			fmt.Println(filePath)
		}
	}
}

func (fileGroup *Group) HasFile() bool {
	return fileGroup.ManifestFile != ""
}

func (fileGroup *Group) HasLockFiles() bool {
	return len(fileGroup.LockFiles) > 0
}

func (fileGroup *Group) GetAllFiles() []string {
	var files []string
	if fileGroup.HasFile() {
		files = append(files, fileGroup.ManifestFile)
	}

	return append(files, fileGroup.LockFiles...)
}

func (fileGroup *Group) matchLockFile(lockFile, dir string) bool {
	if !fileGroup.HasFile() {

		return false
	}
	groupDir, manifestFile := filepath.Split(fileGroup.ManifestFile)
	if groupDir != dir {

		return false
	}

	return matchFile(manifestFile, lockFile)
}

func (fileGroup *Group) matchManifestFile(manifestFile, dir string) bool {
	if fileGroup.HasFile() {

		return false
	}
	groupDir, lockFile := filepath.Split(fileGroup.LockFiles[0])
	if groupDir != dir {

		return false
	}

	return matchFile(manifestFile, lockFile)
}

func matchFile(manifestFile, lockFile string) bool {
	var isPIP bool
	lockFile, isPIP = strings.CutSuffix(lockFile, ".pip.debricked.lock")
	if isPIP {
		lockFile, _ = strings.CutPrefix(lockFile, ".")

		return lockFile == manifestFile
	}

	return true
}
