package testdata

import (
	"os"
	"path/filepath"

	"github.com/debricked/cli/internal/callgraph"
	"github.com/debricked/cli/internal/callgraph/cgexec"
)

type GeneratorMock struct {
	Err         error
	files       []string
	LastOptions callgraph.DebrickedOptions
}

func (r *GeneratorMock) GenerateWithTimer(options callgraph.DebrickedOptions) error {
	r.LastOptions = options

	return r.Err
}

func (r *GeneratorMock) Generate(options callgraph.DebrickedOptions, _ cgexec.IContext) error {
	r.LastOptions = options

	for _, f := range r.files {
		createdFile, err := os.Create(f)
		if err != nil {
			return err
		}

		err = createdFile.Close()
		if err != nil {
			return err
		}
	}

	return r.Err
}

func (r *GeneratorMock) SetFiles(files []string) {
	r.files = files
}

func (r *GeneratorMock) CleanUp() error {
	for _, f := range r.files {
		abs, err := filepath.Abs(f)
		if err != nil {
			return err
		}
		err = os.Remove(abs)
		if err != nil {
			return err
		}
	}

	return nil
}
