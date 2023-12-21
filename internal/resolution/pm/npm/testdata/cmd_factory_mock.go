package testdata

import (
	"os/exec"
)

type CmdFactoryMock struct {
	InstallCmdName string
	MakeInstallErr error
}

func NewEchoCmdFactory() CmdFactoryMock {
	return CmdFactoryMock{
		InstallCmdName: "echo",
	}
}

func (f CmdFactoryMock) MakeInstallCmd(command string, file string) (*exec.Cmd, error) {
	return exec.Command(f.InstallCmdName), f.MakeInstallErr
}
