package testdata

import "os/exec"

type CmdFactoryMock struct {
	MvnCopyDepName   string
	MvnCopyDepErr    error
	CallGraphGenName string
	CallGraphGenErr  error
}

func NewEchoCmdFactory() CmdFactoryMock {
	return CmdFactoryMock{
		MvnCopyDepName:   "echo",
		CallGraphGenName: "echo",
	}
}

func (f CmdFactoryMock) MakeMvnCopyDependenciesCmd(_ string, _ string) (*exec.Cmd, error) {
	return exec.Command(f.MvnCopyDepName, "MvnCopyDep"), f.MvnCopyDepErr
}

func (f CmdFactoryMock) MakeCallGraphGenerationCmd(_ string, _ string, _ string, _ string) (*exec.Cmd, error) {
	return exec.Command(f.CallGraphGenName, "CallGraphGen"), f.CallGraphGenErr
}
