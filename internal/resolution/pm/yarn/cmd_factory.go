package yarn

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
		Args: []string{command, "install",
			"--non-interactive",  // We can't answer any prompts...
			"--ignore-scripts",   // Avoid risky scripts
			"--ignore-engines",   // We won't run the code, so we don't care about the engine versions
			"--ignore-platform",  // We won't run the code, so we don't care about the platform, undocumented option
			"--no-bin-links",     // We don't need symlinks to binaries as we won't run any code
			"--production=false", // Always include dev dependencies
		},
		Dir: fileDir,
	}, err
}
