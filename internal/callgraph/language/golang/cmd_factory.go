package golang

import (
	"os/exec"

	"github.com/debricked/cli/internal/callgraph/cgexec"
)

type ICmdFactory interface {
	MakeCallGraphGenerationCmd(pathToMain string, workingDirectory string, ctx cgexec.IContext) (*exec.Cmd, error)
}

type CmdFactory struct{}

func (_ CmdFactory) MakeCallGraphGenerationCmd(pathToMain string, workingDirectory string, ctx cgexec.IContext) (*exec.Cmd, error) {

	path, err := exec.LookPath("callgraph")

	args := []string{
		"callgraph",
		"-format='{{.Caller}} file:{{.Filename}} {{.Line}} {{.Column}}--->{{.Callee}}'",
		"-algo",
		"cha",
		pathToMain,
	}

	return cgexec.MakeCommand(workingDirectory, path, args, ctx), err
}
