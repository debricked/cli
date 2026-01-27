package testdata

import (
	"os/exec"
	"runtime"
)

type CmdFactoryMock struct {
	Err  error
	Name string
	Arg  string
}

func (f CmdFactoryMock) MakeLockCmd(_ string) (*exec.Cmd, error) {
	if len(f.Arg) == 0 {
		f.Arg = `"MakeLockCmd"`
	}

	if runtime.GOOS == "windows" && f.Name == "echo" {
		return exec.Command("cmd", "/C", f.Name, f.Arg), f.Err
	}

	return exec.Command(f.Name, f.Arg), f.Err
}
