package testdata

import "os/exec"

type CmdFactoryMock struct {
	ListCmdName    string
	MakeListCmdErr error
	ShowCmdName    string
	MakeShowCmdErr error
}

func (f CmdFactoryMock) MakeListCmd() (*exec.Cmd, error) {
	return exec.Command(f.ListCmdName, `"MakeListCmd"`), f.MakeListCmdErr
}

func (f CmdFactoryMock) MakeListCmd() (*exec.Cmd, error) {
	return exec.Command(f.ShowCmdName, `"MakeShowCmd"`), f.MakeShowCmdErr
}
