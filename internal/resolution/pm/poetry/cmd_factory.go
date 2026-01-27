package poetry

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
	poetryPath, err := cmdf.execPath.LookPath("poetry")
	if err != nil {
		return nil, err
	}

	workingDir := filepath.Dir(filepath.Clean(manifestFile))

	env := os.Environ()
	env = append(env, "POETRY_VIRTUALENVS_CREATE=false")

	return &exec.Cmd{
		Path: poetryPath,
		Args: []string{"poetry", "lock"},
		Dir:  workingDir,
		Env:  env,
	}, nil
}
