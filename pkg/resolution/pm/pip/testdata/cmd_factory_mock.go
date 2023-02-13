package testdata

import "os/exec"

type CmdFactoryMock struct {
	GraphCmdName    string
	MakeGraphCmdErr error
	ListCmdName     string
	MakeListCmdErr  error
}

func (f CmdFactoryMock) MakeGraphCmd() (*exec.Cmd, error) {
	return exec.Command(f.GraphCmdName, `"MakeGraphCmd"`), f.MakeGraphCmdErr
}

func (f CmdFactoryMock) MakeListCmd() (*exec.Cmd, error) {
	return exec.Command(f.ListCmdName, `"MakeListCmd"`), f.MakeListCmdErr
}
