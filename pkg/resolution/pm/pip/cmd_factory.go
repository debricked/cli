package pip

import (
	"os/exec"
	"strings"
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

	return &exec.Cmd{
		Path: path,
		Args: []string{"pip", "show", strings.Join(list, " ")},
	}, err
}
