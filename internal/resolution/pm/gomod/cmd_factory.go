package gomod

import "os/exec"

type ICmdFactory interface {
	MakeGraphCmd(workingDirectory string) (*exec.Cmd, error)
	MakeListCmd(workingDirectory string) (*exec.Cmd, error)
	MakeListJsonCmd(workingDirectory string) (*exec.Cmd, error)
}

type CmdFactory struct{}

func (_ CmdFactory) MakeGraphCmd(workingDirectory string) (*exec.Cmd, error) {
	path, err := exec.LookPath("go")

	return &exec.Cmd{
		Path: path,
		Args: []string{"go", "mod", "graph"},
		Dir:  workingDirectory,
	}, err
}

func (_ CmdFactory) MakeListCmd(workingDirectory string) (*exec.Cmd, error) {
	path, err := exec.LookPath("go")

	return &exec.Cmd{
		Path: path,
		Args: []string{"go", "list", "-mod=readonly", "-e", "-m", "all"},
		Dir:  workingDirectory,
	}, err
}

func (_ CmdFactory) MakeListJsonCmd(workingDirectory string) (*exec.Cmd, error) {
	path, err := exec.LookPath("go")

	return &exec.Cmd{
		Path: path,
		Args: []string{path, "list", "-json", "./..."},
		Dir:  workingDirectory,
	}, err
}
