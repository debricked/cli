package uv

import (
	"os"
	"os/exec"
	"path/filepath"
)

type ICmdFactory interface {
	MakeLockCmd(manifestFile string) (*exec.Cmd, error)
}

type IExecPath interface {
	LookPath(file string) (string, error)
}

type ExecPath struct{}

func (_ ExecPath) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

type CmdFactory struct {
	execPath IExecPath
}

func (cmdf CmdFactory) MakeLockCmd(manifestFile string) (*exec.Cmd, error) {
	uvPath, err := cmdf.execPath.LookPath("uv")
	if err != nil {
		return nil, err
	}

	workingDir := filepath.Dir(filepath.Clean(manifestFile))

	env := os.Environ()

	return &exec.Cmd{
		Path: uvPath,
		Args: []string{"uv", "lock"},
		Dir:  workingDir,
		Env:  env,
	}, nil
}
