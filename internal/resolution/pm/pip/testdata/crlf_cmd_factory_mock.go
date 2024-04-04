package testdata

import (
	"os/exec"
)

type CmdFactoryMockCRLF struct {
	CreateVenvCmdName string
	MakeCreateVenvErr error
	InstallCmdName    string
	MakeInstallErr    error
	CatCmdName        string
	MakeCatErr        error
	ListCmdName       string
	MakeListErr       error
	ShowCmdName       string
	MakeShowErr       error
}

func NewCRLFEchoCmdFactory() CmdFactoryMockCRLF {
	return CmdFactoryMockCRLF{
		CreateVenvCmdName: "echo",
		InstallCmdName:    "echo",
		CatCmdName:        "echo",
		ListCmdName:       "echo",
		ShowCmdName:       "echo",
	}
}

func (f CmdFactoryMockCRLF) MakeCreateVenvCmd(file string) (*exec.Cmd, error) {
	return exec.Command(f.CreateVenvCmdName, file), f.MakeCreateVenvErr
}

func (f CmdFactoryMockCRLF) MakeInstallCmd(command string, file string) (*exec.Cmd, error) {
	return exec.Command(f.InstallCmdName, file), f.MakeInstallErr
}

func (f CmdFactoryMockCRLF) MakeListCmd(command string) (*exec.Cmd, error) {
	return exec.Command(f.ListCmdName, "Package                       Version      Editable project location\r\n----------------------------- ------------ ------------------------------------------------------\r\nFlask                         2.0.3"), f.MakeListErr
}

func (f CmdFactoryMockCRLF) MakeCatCmd(file string) (*exec.Cmd, error) {
	return exec.Command(f.CatCmdName, "Flask==2.1.5\r\n"), f.MakeCatErr
}

func (f CmdFactoryMockCRLF) MakeShowCmd(command string, list []string) (*exec.Cmd, error) {
	return exec.Command(f.ShowCmdName, "Name: Flask\r\nVersion: 2.1.2\r\nSummary: A simple framework for building complex web applications.\r\nHome-page: https://palletsprojects.com/p/flask\r\nAuthor: Armin Ronacher\r\nAuthor-email: armin.ronacher@active-4.com\r\nLicense: BSD-3-Clause\r\nLocation: /path/to/site-packages\r\nRequires: click, importlib-metadata, itsdangerous, Jinja2, Werkzeug\r\nRequired-by: Flask-Script, Flask-Compress, Flask-Bcrypt\r\n"), f.MakeShowErr
}
