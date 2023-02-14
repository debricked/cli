package pip

import (
	"os/exec"
)

type ICmdFactory interface {
	MakeListCmd() (*exec.Cmd, error)
	MakeShowCmd(list []string) (*exec.Cmd, error)
}

type CmdFactory struct{}

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
		//expand list with ...
		Args: args,
	}, err
}
