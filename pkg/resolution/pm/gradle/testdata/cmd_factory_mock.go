package testdata

import "os/exec"

type CmdFactoryMock struct {
	Err  error
	Name string
}

func (f CmdFactoryMock) MakeDependenciesGraphCmd(_ string, _ string, _ string) (*exec.Cmd, error) {
	return exec.Command(f.Name, `MakeDependenciesCmd`), f.Err
}

// implement the interface
func (f CmdFactoryMock) MakeFindSubGraphCmd(_ string, _ string, _ string) (*exec.Cmd, error) {
	return exec.Command(f.Name, `MakeFindSubGraphCmd`), f.Err
}

// implement the interface
func (f CmdFactoryMock) MakeDependenciesCmd(_ string) (*exec.Cmd, error) {
	return exec.Command(f.Name, `MakeDependenciesCmd`), f.Err
}
