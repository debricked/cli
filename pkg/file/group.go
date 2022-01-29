package file

import (
	"fmt"
)

type Group struct {
	FilePath     string
	Format       *Format
	RelatedFiles []string
}

func (fileGroup *Group) Print() {
	fmt.Println(fileGroup.FilePath)

	for _, filePath := range fileGroup.RelatedFiles {
		fmt.Println(" * " + filePath)
	}
}

func NewGroup(filePath string, format *Format, relatedFiles []string) *Group {
	return &Group{FilePath: filePath, Format: format, RelatedFiles: relatedFiles}
}
