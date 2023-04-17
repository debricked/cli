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

func (fileGroup *Group) checkFilePathDependantCases(fileMatch bool, lockFileMatch bool, file string) bool {
	if lockFileMatch {
		filePathDependantCases := fileGroup.getFilePathDependantCases()
		for _, c := range filePathDependantCases {
			if strings.HasSuffix(file, c) {
				fileBase, _ := strings.CutSuffix(file, c)

				return len(fileGroup.ManifestFile) > 0 && (fileBase == filepath.Base(fileGroup.ManifestFile))
			}
		}

		return true
	}

	if fileMatch {
		filePathDependantCases := fileGroup.getFilePathDependantCases()
		for _, c := range filePathDependantCases {
			for _, lockFile := range fileGroup.LockFiles {
				if strings.HasSuffix(lockFile, c) {
					lockFileBase, _ := strings.CutSuffix(lockFile, c)

					return lockFileBase == file
				}
			}
		}

		return true
	}

	return true
}

func (fileGroup *Group) getFilePathDependantCases() []string {
	return []string{
		".pip.debricked.lock",
	}
}
