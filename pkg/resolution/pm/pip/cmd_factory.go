package pip

import (
	"os/exec"
)

type ICmdFactory interface {
	MakeInstallCmd(file string) (*exec.Cmd, error)
	MakeCatCmd(file string) (*exec.Cmd, error)
	MakeListCmd() (*exec.Cmd, error)
	MakeShowCmd(list []string) (*exec.Cmd, error)
}

type CmdFactory struct{}

func (_ CmdFactory) MakeCatCmd(file string) (*exec.Cmd, error) {
	path, err := exec.LookPath("cat")

	return &exec.Cmd{
		Path: path,
		Args: []string{"cat", file},
	}, err
}

func (_ CmdFactory) MakeListCmd() (*exec.Cmd, error) {
	path, err := exec.LookPath("pip")

	return &exec.Cmd{
		Path: path,
		Args: []string{"pip", "list"},
	}, err
}

func (_ CmdFactory) MakeShowCmd(list []string) (*exec.Cmd, error) {
	path, err := exec.LookPath("pip")

	args := []string{"pip", "show"}
	args = append(args, list...)

	return &exec.Cmd{
		Path: path,
		Args: args,
	}, err
}

func (_ CmdFactory) MakeInstallCmd(file string) (*exec.Cmd, error) {
	path, err := exec.LookPath("pip")

	return &exec.Cmd{
		Path: path,
		Args: []string{"pip", "install", "-r", file},
	}, err
}
