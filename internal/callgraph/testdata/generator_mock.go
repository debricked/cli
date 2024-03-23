package testdata

import (
	"os"
	"path/filepath"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	"github.com/debricked/cli/internal/callgraph/config"
)

type GeneratorMock struct {
	Err   error
	files []string
}

func (r *GeneratorMock) GenerateWithTimer(_ []string, _ []string, _ []string, _ []config.IConfig, _ int) error {
	return r.Err
}

func (r *GeneratorMock) Generate(_ []string, _ []string, _ []string, _ []config.IConfig, _ cgexec.IContext) error {
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
