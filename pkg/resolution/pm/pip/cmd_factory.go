package pip

import (
	"fmt"
	"os/exec"
)

type ICmdFactory interface {
	MakeCreateVenvCmd(file string) (*exec.Cmd, error)
	MakeInstallCmd(command string, file string) (*exec.Cmd, error)
	MakeCatCmd(file string) (*exec.Cmd, error)
	MakeListCmd(command string) (*exec.Cmd, error)
	MakeShowCmd(command string, list []string) (*exec.Cmd, error)
}

type CmdFactory struct{}

func (_ CmdFactory) MakeCreateVenvCmd(fpath string) (*exec.Cmd, error) {
	path, err := exec.LookPath("python")

	return &exec.Cmd{
		Path: path,
		Args: []string{"python", "-m", "venv", fpath, "--clear"},
	}, err
}

func (_ CmdFactory) MakeInstallCmd(command string, file string) (*exec.Cmd, error) {
	path, err := exec.LookPath(command)

	return &exec.Cmd{
		Path: path,
		Args: []string{command, "install", "-r", file},
	}, err
}

func (_ CmdFactory) MakeCatCmd(file string) (*exec.Cmd, error) {
	path, err := exec.LookPath("cat")

	return &exec.Cmd{
		Path: path,
		Args: []string{"cat", file},
	}, err
}

func (_ CmdFactory) MakeListCmd(command string) (*exec.Cmd, error) {
	path, err := exec.LookPath(command)
	fmt.Println("MakeListCmd", command)

	return &exec.Cmd{
		Path: path,
		Args: []string{"pip", "list"},
	}, err
}

func (_ CmdFactory) MakeShowCmd(command string, list []string) (*exec.Cmd, error) {
	path, err := exec.LookPath(command)
	fmt.Println("MakeShowCmd", command)

	args := []string{command, "show"}
	args = append(args, list...)

	return &exec.Cmd{
		Path: path,
		Args: args,
	}, err
}
