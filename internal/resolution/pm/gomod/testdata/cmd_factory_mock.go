package testdata

import (
	"os/exec"
	"strings"
)

type CmdFactoryMock struct {
	GraphCmdName       string
	MakeGraphCmdErr    error
	ListCmdName        string
	MakeListCmdErr     error
	ListJsonCmdName    string
	MakeListJsonCmdErr error
}

func NewEchoCmdFactory() CmdFactoryMock {
	return CmdFactoryMock{
		GraphCmdName:    "echo",
		ListCmdName:     "echo",
		ListJsonCmdName: "echo",
	}
}

func (f CmdFactoryMock) MakeGraphCmd(_ string) (*exec.Cmd, error) {
	return exec.Command(f.GraphCmdName, "MakeGraphCmd"), f.MakeGraphCmdErr
}

func (f CmdFactoryMock) MakeListCmd(_ string) (*exec.Cmd, error) {
	output := strings.Join(
		[]string{
			"github.com/debricked/cli",
			"module1 v1.1.0",
			"module2 v0.4.0 => ./localDepOne",
			"module3 v1.1.4",
			"module4 v0.0.0-20170915032832-14c0d48ead0c",
		}, "\n")

	return exec.Command(f.ListCmdName, output), f.MakeListCmdErr
}

func (f CmdFactoryMock) MakeListJsonCmd(_ string) (*exec.Cmd, error) {
	output := strings.Join(
		[]string{
			`{"ImportPath": "module1/package1", "TestImports": ["module4/package3"], "Imports": ["module2/package3", "fmt", "sync"]}`,
			`{"ImportPath": "module2/package2"}`,
			`{"ImportPath": "module3/package3", "TestImports": ["module4/package1", "module2/package3", "os", "path/filepath", "testing"]}`,
			`{"ImportPath": "module4/package4", "TestImports": []}`,
			`{"ImportPath": "missing/package", "TestImports": [], "Imports": ["sync"]}`,
		}, "\n")

	return exec.Command(f.ListJsonCmdName, output), f.MakeListJsonCmdErr
}
