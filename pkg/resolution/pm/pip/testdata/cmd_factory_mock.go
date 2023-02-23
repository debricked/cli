package testdata

import (
	"os"
	"os/exec"
)

type CmdFactoryMock struct {
	CreateVenvCmdName   string
	MakeCreateVenvErr   error
	ActivateVenvCmdName string
	MakeActivateVenvErr error
	InstallCmdName      string
	MakeInstallErr      error
	CatCmdName          string
	MakeCatErr          error
	ListCmdName         string
	MakeListErr         error
	ShowCmdName         string
	MakeShowErr         error
}

func NewEchoCmdFactory() CmdFactoryMock {
	return CmdFactoryMock{
		CreateVenvCmdName:   "echo",
		ActivateVenvCmdName: "echo",
		InstallCmdName:      "echo",
		CatCmdName:          "echo",
		ListCmdName:         "echo",
		ShowCmdName:         "echo",
	}
}

func (f CmdFactoryMock) MakeCreateVenvCmd(file string) (*exec.Cmd, error) {
	return exec.Command(f.CreateVenvCmdName, file), f.MakeCreateVenvErr
}

func (f CmdFactoryMock) MakeActivateVenvCmd(file string) (*exec.Cmd, error) {
	return exec.Command(f.ActivateVenvCmdName, file), f.MakeActivateVenvErr
}

func (f CmdFactoryMock) MakeInstallCmd(command string, file string) (*exec.Cmd, error) {
	return exec.Command(f.InstallCmdName, file), f.MakeInstallErr
}

func (f CmdFactoryMock) MakeListCmd(command string) (*exec.Cmd, error) {
	fileContent, err := os.ReadFile("testdata/list.txt")
	if err != nil {
		return nil, err
	}
	pipData := string(fileContent)
	return exec.Command(f.ListCmdName, pipData), f.MakeListErr
}

func (f CmdFactoryMock) MakeCatCmd(file string) (*exec.Cmd, error) {
	fileContent, err := os.ReadFile("testdata/requirements.txt")
	if err != nil {
		return nil, err
	}
	requirements := string(fileContent)
	return exec.Command(f.CatCmdName, requirements), f.MakeCatErr
}

func (f CmdFactoryMock) MakeShowCmd(command string, list []string) (*exec.Cmd, error) {
	fileContent, err := os.ReadFile("testdata/show.txt")
	if err != nil {
		return nil, err
	}
	show := string(fileContent)
	return exec.Command(f.ShowCmdName, show), f.MakeShowErr
}
