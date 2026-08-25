package testdata

import (
	"os/exec"
	"path/filepath"
	"runtime"
)

type CmdFactoryMock struct {
	LockErr error
	DepsErr error
	Name    string
	Arg     string
	// DepsFile makes MakeDepsCmd print the contents of the given file, which
	// allows tests to feed a realistic dependency tree to the job.
	DepsFile string
}

func (f CmdFactoryMock) MakeResolveCmd(_ string) (*exec.Cmd, error) {
	if len(f.Arg) == 0 {
		f.Arg = `"MakeResolveCmd"`
	}

	if runtime.GOOS == "windows" && f.Name == "echo" {
		return exec.Command("cmd", "/C", f.Name, f.Arg), f.LockErr
	}

	return exec.Command(f.Name, f.Arg), f.LockErr
}

func (f CmdFactoryMock) MakeDepsCmd(_ string) (*exec.Cmd, error) {
	if len(f.DepsFile) > 0 {
		if runtime.GOOS == "windows" {
			return exec.Command("cmd", "/C", "type", filepath.FromSlash(f.DepsFile)), f.DepsErr
		}

		return exec.Command("cat", f.DepsFile), f.DepsErr
	}

	if len(f.Arg) == 0 {
		f.Arg = `"MakeDepsCmd"`
	}

	if runtime.GOOS == "windows" && f.Name == "echo" {
		return exec.Command("cmd", "/C", f.Name, f.Arg), f.DepsErr
	}

	return exec.Command(f.Name, f.Arg), f.DepsErr
}
