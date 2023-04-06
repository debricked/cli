package file

import (
	"fmt"
	"path/filepath"
	"strings"
)

type Group struct {
	FilePath       string          `json:"dependencyFile"`
	CompiledFormat *CompiledFormat `json:"-"`
	RelatedFiles   []string        `json:"lockFiles"`
}

func NewGroup(filePath string, format *CompiledFormat, relatedFiles []string) *Group {
	return &Group{FilePath: filePath, CompiledFormat: format, RelatedFiles: relatedFiles}
}

func (fileGroup *Group) Print() {
	hasFile := fileGroup.HasFile()
	if hasFile {
		fmt.Println(fileGroup.FilePath)
	}
	for _, filePath := range fileGroup.RelatedFiles {
		if hasFile {
			fmt.Println(" * " + filePath)
		} else {
			fmt.Println(filePath)
		}
	}
}

func (fileGroup *Group) HasFile() bool {
	return fileGroup.FilePath != ""
}

func (fileGroup *Group) HasLockFiles() bool {
	return len(fileGroup.RelatedFiles) > 0
}

func (fileGroup *Group) GetAllFiles() []string {
	var files []string
	if fileGroup.HasFile() {
		files = append(files, fileGroup.FilePath)
	}

	return append(files, fileGroup.RelatedFiles...)
}

func (fileGroup *Group) checkFilePathDependantCases(fileMatch bool, lockFileMatch bool, file string) bool {
	if lockFileMatch {
		filePathDependantCases := fileGroup.getFilePathDependantCases()
		for _, c := range filePathDependantCases {
			if strings.HasSuffix(file, c) {
				fileBase, _ := strings.CutSuffix(file, c)

				return len(fileGroup.FilePath) > 0 && (fileBase == filepath.Base(fileGroup.FilePath))
			}
		}

		return true
	}

	if fileMatch {
		filePathDependantCases := fileGroup.getFilePathDependantCases()
		for _, c := range filePathDependantCases {
			for _, lockFile := range fileGroup.RelatedFiles {
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
