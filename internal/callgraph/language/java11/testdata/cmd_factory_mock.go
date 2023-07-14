package testdata

import (
	"os/exec"

	"github.com/debricked/cli/internal/callgraph/cgexec"
)

type CmdFactoryMock struct {
	MvnCopyDepName   string
	MvnCopyDepErr    error
	CallGraphGenName string
	CallGraphGenErr  error
	BuildMavenName   string
	BuildMavenErr    error
}

func NewEchoCmdFactory() CmdFactoryMock {
	return CmdFactoryMock{
		MvnCopyDepName:   "echo",
		CallGraphGenName: "echo",
		BuildMavenName:   "echo",
	}
}

func (f CmdFactoryMock) MakeMvnCopyDependenciesCmd(_ string, _ string, _ cgexec.IContext) (*exec.Cmd, error) {
	return exec.Command(f.MvnCopyDepName, "MvnCopyDep"), f.MvnCopyDepErr
}

func (f CmdFactoryMock) MakeCallGraphGenerationCmd(_ string, _ string, _ []string, _ string, _ string, _ cgexec.IContext) (*exec.Cmd, error) {
	return exec.Command(f.CallGraphGenName, "CallGraphGen"), f.CallGraphGenErr
}

func (f CmdFactoryMock) MakeBuildMavenCmd(workingDirectory string, ctx cgexec.IContext) (*exec.Cmd, error) {
	return exec.Command(f.BuildMavenName, "BuildMaven"), f.BuildMavenErr
}
