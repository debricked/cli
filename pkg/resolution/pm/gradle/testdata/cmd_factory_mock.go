package testdata

import "os/exec"

type CmdFactoryMock struct {
	Err  error
	Name string
}

func (f CmdFactoryMock) MakeDependenciesCmd() (*exec.Cmd, error) {
	return exec.Command(f.Name, `MakeDependenciesCmd`), f.Err
}
