package testdata

import "os/exec"

type CmdFactoryMock struct {
	GradleCopyDepName string
	GradleCopyDepErr  error
	MvnCopyDepName    string
	MvnCopyDepErr     error
	CallGraphGenName  string
	CallGraphGenErr   error
}

func NewEchoCmdFactory() CmdFactoryMock {
	return CmdFactoryMock{
		GradleCopyDepName: "echo",
		MvnCopyDepName:    "echo",
		CallGraphGenName:  "echo",
	}
}

func (f CmdFactoryMock) MakeGradleCopyDependenciesCmd(_ string, _ string, _ string) (*exec.Cmd, error) {
	return exec.Command(f.GradleCopyDepName, "GradleCopyDep"), f.GradleCopyDepErr
}

func (f CmdFactoryMock) MakeMvnCopyDependenciesCmd(_ string, _ string) (*exec.Cmd, error) {
	return exec.Command(f.MvnCopyDepName, "MvnCopyDep"), f.MvnCopyDepErr
}

func (f CmdFactoryMock) MakeCallGraphGenerationCmd(_ string, _ string, _ string, _ string) (*exec.Cmd, error) {
	return exec.Command(f.CallGraphGenName, "CallGraphGen"), f.CallGraphGenErr
}
