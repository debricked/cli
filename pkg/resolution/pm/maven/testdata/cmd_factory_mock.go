package testdata

import "os/exec"

type CmdFactoryMock struct {
	Err  error
	Name string
	Arg  string
}

func (f CmdFactoryMock) MakeDependencyTreeCmd(_ string) (*exec.Cmd, error) {
	if len(f.Arg) == 0 {
		f.Arg = `"MakeDependencyTreeCmd"`
	}
	return exec.Command(f.Name, f.Arg), f.Err
}
