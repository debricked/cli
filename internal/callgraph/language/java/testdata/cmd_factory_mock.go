package testdata

import (
	"os/exec"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	ioFs "github.com/debricked/cli/internal/io"
)

type CmdFactoryMock struct {
	MvnCopyDepName   string
	MvnCopyDepErr    error
	CallGraphGenName string
	CallGraphGenErr  error
	BuildMavenName   string
	BuildMavenErr    error
	JavaVersionName  string
	JavaVersionErr   error
}

func NewEchoCmdFactory() CmdFactoryMock {
	return CmdFactoryMock{
		MvnCopyDepName:   "echo",
		CallGraphGenName: "echo",
		BuildMavenName:   "echo",
		JavaVersionName:  "echo",
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

func (f CmdFactoryMock) MakeJavaVersionCmd(workingDirectory string, ctx cgexec.IContext) (*exec.Cmd, error) {
	return exec.Command(
		f.JavaVersionName,
		"\"openjdk 23.0.1 2024-10-15\nOpenJDK Runtime Environment Homebrew (build 23.0.1)\nOpenJDK 64-Bit Server VM Homebrew (build 23.0.1, mixed mode, sharing)\"",
	), f.JavaVersionErr
}

type MockSootHandler struct {
	GetSootWrapperError error
}

func (msh MockSootHandler) GetSootWrapper(version string, fs ioFs.IFileSystem, arc ioFs.IArchive) (string, error) {
	return "", msh.GetSootWrapperError
}
