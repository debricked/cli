package testdata

import (
	"os/exec"
	"strings"
)

type CmdFactoryMock struct {
	Err  error
	Name string
}

func (f CmdFactoryMock) MakeDependenciesGraphCmd(dir string, gradlew string, _ string) (*exec.Cmd, error) {
	err := f.Err
	if gradlew == "gradle" {
		err = nil
	}

	if f.Err != nil && strings.HasPrefix(f.Err.Error(), "give-error-on-gradle") {
		err = f.Err
	}

	return exec.Command(f.Name, `MakeDependenciesCmd`), err
}

// implement the interface
func (f CmdFactoryMock) MakeFindSubGraphCmd(_ string, _ string, _ string) (*exec.Cmd, error) {
	return exec.Command(f.Name, `MakeFindSubGraphCmd`), f.Err
}

// implement the interface
func (f CmdFactoryMock) MakeDependenciesCmd(_ string) (*exec.Cmd, error) {
	return exec.Command(f.Name, `MakeDependenciesCmd`), f.Err
}
