package pub

import (
	"os"
	"os/exec"
	"path/filepath"
)

type ICmdFactory interface {
	MakeLockCmd(manifestFile string) (*exec.Cmd, error)
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

// MakeLockCmd creates an exec.Cmd that runs `dart pub get` in the directory
// of the given manifest file. This resolves all dependencies and writes a
// pubspec.lock file that can be used for SCA scanning.
func (cmdf CmdFactory) MakeLockCmd(manifestFile string) (*exec.Cmd, error) {
	dartPath, err := cmdf.execPath.LookPath("dart")
	if err != nil {
		return nil, err
	}

	workingDir := filepath.Dir(filepath.Clean(manifestFile))

	return &exec.Cmd{
		Path: dartPath,
		Args: []string{"dart", "pub", "get"},
		Dir:  workingDir,
		Env:  os.Environ(),
	}, nil
}

// MakeDepsCmd creates an exec.Cmd that runs `dart pub deps --json` in the
// directory of the given manifest file. This outputs explicit parent-child
// dependency relationships used to reconstruct the transitive dependency tree.
func (cmdf CmdFactory) MakeDepsCmd(manifestFile string) (*exec.Cmd, error) {
	dartPath, err := cmdf.execPath.LookPath("dart")
	if err != nil {
		return nil, err
	}

	workingDir := filepath.Dir(filepath.Clean(manifestFile))

	return &exec.Cmd{

		Path: dartPath,
		Args: []string{"dart", "pub", "deps", "--json"},
		Dir:  workingDir,
		Env:  os.Environ(),
	}, nil
}
