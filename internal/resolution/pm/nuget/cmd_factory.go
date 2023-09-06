package nuget

import (
	"os/exec"
	"path/filepath"
)

type ICmdFactory interface {
	MakeInstallCmd(command string, file string) (*exec.Cmd, error)
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
		Args: []string{command, "restore",
			"--use-lock-file",
		},
		Dir: fileDir,
	}, err
}
