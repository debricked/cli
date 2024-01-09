package testdata

import (
	"os/exec"
)

type EmptyCmdFactoryMock struct {
	MakeErr error
}

func NewEmptyCmdFactory() EmptyCmdFactoryMock {
	return EmptyCmdFactoryMock{}
}

func (f EmptyCmdFactoryMock) MakeInstallCmd(_ string, _ string) (*exec.Cmd, error) {
	return nil, f.MakeErr
}

func (f EmptyCmdFactoryMock) GetTempoCsproj() string {
	return ""
}
