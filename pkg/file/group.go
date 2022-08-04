package file

import (
	"fmt"
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

func (fileGroup *Group) GetAllFiles() []string {
	var files []string
	if fileGroup.HasFile() {
		files = append(files, fileGroup.FilePath)
	}

	return append(files, fileGroup.RelatedFiles...)
}
