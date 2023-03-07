package gomod

import "os/exec"

type ICmdFactory interface {
	MakeTidyCmd() (*exec.Cmd, error)
	MakeGraphCmd() (*exec.Cmd, error)
	MakeListCmd() (*exec.Cmd, error)
}

type CmdFactory struct{}

func (_ CmdFactory) MakeTidyCmd() (*exec.Cmd, error) {
	path, err := exec.LookPath("go")

	return &exec.Cmd{
		Path: path,
		Args: []string{"go", "mod", "tidy"},
	}, err
}

func (_ CmdFactory) MakeGraphCmd() (*exec.Cmd, error) {
	path, err := exec.LookPath("go")

	return &exec.Cmd{
		Path: path,
		Args: []string{"go", "mod", "graph"},
	}, err
}

func (_ CmdFactory) MakeListCmd() (*exec.Cmd, error) {
	path, err := exec.LookPath("go")

	return &exec.Cmd{
		Path: path,
		Args: []string{"go", "list", "-mod=readonly", "-e", "-m", "all"},
	}, err
}
