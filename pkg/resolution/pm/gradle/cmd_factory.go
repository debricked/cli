package gradle

import "os/exec"

type ICmdFactory interface {
	MakeDependenciesCmd() (*exec.Cmd, error)
}

type CmdFactory struct{}

func (_ CmdFactory) MakeDependenciesCmd() (*exec.Cmd, error) {
	path, err := exec.LookPath("gradle")

	return &exec.Cmd{
		Path: path,
		Args: []string{"gradle", "dependencies"},
	}, err
}
