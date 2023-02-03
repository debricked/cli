package testdata

import "os/exec"

type CmdFactoryMock struct {
	Err  error
	Name string
}

func (f CmdFactoryMock) MakeDependencyTreeCmd() (*exec.Cmd, error) {
	return exec.Command(f.Name, `"MakeDependencyTreeCmd"`), f.Err
}
