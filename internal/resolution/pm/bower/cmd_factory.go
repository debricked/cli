package bower

import (
	"os/exec"
	"path/filepath"
)

type ICmdFactory interface {
	MakeInstallCmd(command string, file string) (*exec.Cmd, error)
	MakeListCmd(command string, file string) (*exec.Cmd, error)
}

type IExecPath interface {
	LookPath(file string) (string, error)
}

type ExecPath struct {
}

func (ExecPath) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

type CmdFactory struct {
	execPath IExecPath
}

func (cmdf CmdFactory) MakeInstallCmd(command string, file string) (*exec.Cmd, error) {
	path, err := cmdf.execPath.LookPath(command)

	fileDir := filepath.Dir(file)

	return &exec.Cmd{
		Path: path,
		Args: []string{
			command,
			"install",
			"--save",
			"--save-dev",
			"--save-exact",
			"--allow-root",
		},
		Dir: fileDir,
	}, err
}

func (cmdf CmdFactory) MakeListCmd(command string, file string) (*exec.Cmd, error) {
	path, err := cmdf.execPath.LookPath(command)

	fileDir := filepath.Dir(file)

	return &exec.Cmd{
		Path: path,
		Args: []string{command, "list"},
		Dir:  fileDir,
	}, err
}
