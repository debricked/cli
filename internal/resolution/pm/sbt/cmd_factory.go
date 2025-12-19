package sbt

import "os/exec"

type ICmdFactory interface {
	MakePomCmd(workingDirectory string) (*exec.Cmd, error)
}

type CmdFactory struct{}

func (CmdFactory) MakePomCmd(workingDirectory string) (*exec.Cmd, error) {
	path, err := exec.LookPath("sbt")

	return &exec.Cmd{
		Path: path,
		Args: []string{
			"sbt",
			"makePom",
		},
		Dir: workingDirectory,
	}, err
}
