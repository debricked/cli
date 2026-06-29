package testdata

import (
	"os/exec"
	"runtime"
)

type CmdFactoryMock struct {
	LockErr error
	DepsErr error
	Name    string
	Arg     string
}

func (f CmdFactoryMock) MakeLockCmd(_ string) (*exec.Cmd, error) {
	if len(f.Arg) == 0 {
		f.Arg = `"MakeLockCmd"`
	}

	if runtime.GOOS == "windows" && f.Name == "echo" {
		return exec.Command("cmd", "/C", f.Name, f.Arg), f.LockErr
	}

	return exec.Command(f.Name, f.Arg), f.LockErr
}

func (f CmdFactoryMock) MakeDepsCmd(_ string) (*exec.Cmd, error) {
	if len(f.Arg) == 0 {
		f.Arg = `"MakeDepsCmd"`
	}

	if runtime.GOOS == "windows" && f.Name == "echo" {
		return exec.Command("cmd", "/C", f.Name, f.Arg), f.DepsErr
	}

	return exec.Command(f.Name, f.Arg), f.DepsErr
}
