package testdata

import (
	"os/exec"

	"github.com/debricked/cli/internal/callgraph/cgexec"
)

type CmdFactoryMock struct {
	CallGraphGenName string
	CallGraphGenErr  error
	CommandOutput    string
}

func NewEchoCmdFactory() CmdFactoryMock {
	return CmdFactoryMock{
		CallGraphGenName: "echo",
		CommandOutput:    "CallGraphGen",
	}
}

func (f CmdFactoryMock) MakeCallGraphGenerationCmd(pathToMain string, workingDirectory string, ctx cgexec.IContext) (*exec.Cmd, error) {
	return exec.Command(f.CallGraphGenName, f.CommandOutput), f.CallGraphGenErr
}
