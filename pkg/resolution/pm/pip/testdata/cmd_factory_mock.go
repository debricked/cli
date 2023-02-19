package testdata

import (
	"os"
	"os/exec"
)

type CmdFactoryMock struct {
	InstallCmdName string
	MakeInstallErr error
	CatCmdName     string
	MakeCatErr     error
	ListCmdName    string
	MakeListErr    error
	ShowCmdName    string
	MakeShowErr    error
}

func (f CmdFactoryMock) MakeInstallCmd(file string) (*exec.Cmd, error) {
	return exec.Command(f.InstallCmdName, file), f.MakeListErr
}

func (f CmdFactoryMock) MakeListCmd() (*exec.Cmd, error) {
	fileContent, err := os.ReadFile("testdata/list.txt")
	if err != nil {
		return nil, err
	}
	pipData := `"` + string(fileContent) + `"`
	return exec.Command(f.ListCmdName, pipData), f.MakeListErr
}

func (f CmdFactoryMock) MakeCatCmd(file string) (*exec.Cmd, error) {
	fileContent, err := os.ReadFile("testdata/requirements.txt")
	if err != nil {
		return nil, err
	}
	requirements := `"` + string(fileContent) + `"`
	return exec.Command(f.CatCmdName, requirements), f.MakeCatErr
}

func (f CmdFactoryMock) MakeShowCmd(list []string) (*exec.Cmd, error) {
	fileContent, err := os.ReadFile("testdata/show.txt")
	if err != nil {
		return nil, err
	}
	show := `"` + string(fileContent) + `"`
	return exec.Command(f.ShowCmdName, show), f.MakeShowErr
}
