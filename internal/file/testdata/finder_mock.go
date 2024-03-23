package testdata

import (
	"github.com/debricked/cli/internal/file"
)

type FinderMock struct {
	groups          file.Groups
	compiledFormats []*file.CompiledFormat
	error           error
}

func NewFinderMock() *FinderMock {
	return &FinderMock{
		groups:          file.Groups{},
		compiledFormats: nil,
		error:           nil,
	}
}

// GetGroups return all file groups in specified path recursively.
func (f *FinderMock) GetGroups(_ string, _ []string, _ []string, _ bool, _ int) (file.Groups, error) {
	return f.groups, f.error
}

func (f *FinderMock) GetConfigPath(_ string, _ []string, _ []string) string {
	return ""
}

func (f *FinderMock) GetSupportedFormats() ([]*file.CompiledFormat, error) {
	return f.compiledFormats, f.error
}

func (f *FinderMock) SetGetGroupsReturnMock(gs file.Groups, err error) {
	f.groups = gs
	f.error = err
}

func (f *FinderMock) SetGetSupportedFormatsReturnMock(compiledFormats []*file.CompiledFormat, err error) {
	f.compiledFormats = compiledFormats
	f.error = err
}
