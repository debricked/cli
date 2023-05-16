package file

import "github.com/debricked/cli/internal/resolution/pm"

type IBatch interface {
	Files() []string
	Add(file string)
	Pm() pm.IPm
}

type Batch struct {
	files map[string]bool
	pm    pm.IPm
}

func NewBatch(pm pm.IPm) Batch {
	return Batch{files: make(map[string]bool), pm: pm}
}

func (b Batch) Files() []string {
	var files []string
	for file := range b.files {
		files = append(files, file)
	}

	return files
}

func (b Batch) Add(file string) {
	if ok := b.files[file]; !ok {
		b.files[file] = true
	}
}

func (b Batch) Pm() pm.IPm {
	return b.pm
}
