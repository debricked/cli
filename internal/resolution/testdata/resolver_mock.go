package testdata

import (
	"github.com/debricked/cli/internal/resolution"
	"github.com/debricked/cli/internal/resolution/job"
	"os"
	"path/filepath"
)

type ResolverMock struct {
	Err   error
	files []string
}

func (r *ResolverMock) Resolve(_ []string, _ []string) (resolution.IResolution, error) {
	for _, f := range r.files {
		createdFile, err := os.Create(f)
		if err != nil {
			return nil, err
		}

		err = createdFile.Close()
		if err != nil {
			return nil, err
		}
	}

	return resolution.NewResolution([]job.IJob{}), r.Err
}

func (r *ResolverMock) SetFiles(files []string) {
	r.files = files
}

func (r *ResolverMock) CleanUp() error {
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
