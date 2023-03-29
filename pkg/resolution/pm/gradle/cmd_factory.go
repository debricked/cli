package gradle

import (
	"os/exec"
)

type ICmdFactory interface {
	MakeFindSubGraphCmd(workingDirectory string, gradlew string, initScript string) (*exec.Cmd, error)
	MakeDependenciesGraphCmd(workingDirectory string, gradlew string, initScript string) (*exec.Cmd, error)
}

type CmdFactory struct {
}

func (cf CmdFactory) MakeFindSubGraphCmd(workingDirectory string, gradlew string, initScript string) (*exec.Cmd, error) {
	path, err := exec.LookPath(gradlew)

	return &exec.Cmd{
		Path: path,
		Args: []string{gradlew, "--init-script", initScript, "debrickedFindSubProjectPaths"},
		Dir:  workingDirectory,
	}, err
}

func (cf CmdFactory) MakeDependenciesGraphCmd(workingDirectory string, gradlew string, initScript string) (*exec.Cmd, error) {
	path, err := exec.LookPath(gradlew)

	return &exec.Cmd{
		Path: path,
		Args: []string{gradlew, "--init-script", initScript, "debrickedAllDeps"},
		Dir:  workingDirectory,
	}, err
}
