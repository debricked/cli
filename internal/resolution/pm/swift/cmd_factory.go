package swift

import (
	"os"
	"os/exec"
	"path/filepath"
)

const (
	resolveCmd = "resolve"
	showDeps   = "show-dependencies"
)

type ICmdFactory interface {
	MakeResolveCmd(manifestFile string) (*exec.Cmd, error)
	MakeDepsCmd(manifestFile string) (*exec.Cmd, error)
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

func (cmdf CmdFactory) MakeResolveCmd(manifestFile string) (*exec.Cmd, error) {
	swiftPath, err := cmdf.execPath.LookPath("swift")
	if err != nil {
		return nil, err
	}

	workingDir := filepath.Dir(filepath.Clean(manifestFile))

	return &exec.Cmd{
		Path: swiftPath,
		Args: []string{"swift", "package", resolveCmd},
		Dir:  workingDir,
		Env:  os.Environ(),
	}, nil
}

func (cmdf CmdFactory) MakeDepsCmd(manifestFile string) (*exec.Cmd, error) {
	swiftPath, err := cmdf.execPath.LookPath("swift")
	if err != nil {
		return nil, err
	}

	workingDir := filepath.Dir(filepath.Clean(manifestFile))

	return &exec.Cmd{
		Path: swiftPath,
		Args: []string{"swift", "package", showDeps, "--format", "json"},
		Dir:  workingDir,
		Env:  os.Environ(),
	}, nil
}
