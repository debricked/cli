package file

import (
	"fmt"
)

type Group struct {
	FilePath       string
	CompiledFormat *CompiledFormat
	RelatedFiles   []string
}

func (fileGroup *Group) Print() {
	fmt.Println(fileGroup.FilePath)

	for _, filePath := range fileGroup.RelatedFiles {
		fmt.Println(" * " + filePath)
	}
}

func NewGroup(filePath string, format *CompiledFormat, relatedFiles []string) *Group {
	return &Group{FilePath: filePath, CompiledFormat: format, RelatedFiles: relatedFiles}
}
