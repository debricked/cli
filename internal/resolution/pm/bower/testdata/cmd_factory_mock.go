package testdata

import "os/exec"

type CmdFactoryMock struct {
	InstallCmdName    string
	MakeInstallCmdErr error
	ListCmdName       string
	MakeListCmdErr    error
}

func NewEchoCmdFactory() CmdFactoryMock {
	return CmdFactoryMock{
		InstallCmdName: "echo",
		ListCmdName:    "echo",
	}
}

func (f CmdFactoryMock) MakeInstallCmd(_ string, _ string) (*exec.Cmd, error) {
	return exec.Command(f.InstallCmdName, "MakeInstallCmd"), f.MakeInstallCmdErr
}

func (f CmdFactoryMock) MakeListCmd(_ string, _ string) (*exec.Cmd, error) {
	return exec.Command(f.ListCmdName, "MakeListCmd"), f.MakeListCmdErr
}
