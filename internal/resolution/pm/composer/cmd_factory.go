package composer

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
		Args: []string{command, "update",
			"--no-interaction",       // We can't answer any prompts...
			"--no-scripts",           // Avoid risky scripts
			"--ignore-platform-reqs", // We won't run the code, so we don't care about the platform
			"--no-autoloader",        // We won't execute any code, no need for autoloader
			"--no-install",           // No need to install packages
			"--no-plugins",           // We won't run the code, so no plugins needed
			"--no-audit",             // We don't want to run an audit
		},
		Dir: fileDir,
	}, err
}
