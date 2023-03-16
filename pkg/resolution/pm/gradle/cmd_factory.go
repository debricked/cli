package gradle

import (
	"fmt"
	"os/exec"
)

const initScript = ".debricked.gradle.initscript"

type ICmdFactory interface {
	MakeDependenciesCmd(workingDirectory string) (*exec.Cmd, error)
	MakeFindSubGraphCmd(workingDirectory string) (*exec.Cmd, error)
	MakeDependenciesGraphCmd(workingDirectory string) (*exec.Cmd, error)
}

type CmdFactory struct {
	gradlew    string
	initScript string
}

func (_ CmdFactory) MakeDependenciesCmd(workingDirectory string) (*exec.Cmd, error) {
	path, err := exec.LookPath("gradle")

	return &exec.Cmd{
		Path: path,
		Args: []string{"gradle", "dependencies"},
		Dir:  workingDirectory,
	}, err
}

func (cf CmdFactory) MakeFindSubGraphCmd(workingDirectory string) (*exec.Cmd, error) {
	path, err := exec.LookPath(cf.gradlew)
	fmt.Println(path)
	fmt.Println(err)
	fmt.Println(cf.gradlew, "--init-script", cf.initScript, "debrickedFindSubProjectPaths")

	return &exec.Cmd{
		Path: path,
		Args: []string{cf.gradlew, "--init-script", cf.initScript, "debrickedFindSubProjectPaths"},
		Dir:  workingDirectory,
	}, err
}

func (cf CmdFactory) MakeDependenciesGraphCmd(workingDirectory string) (*exec.Cmd, error) {
	path, err := exec.LookPath(cf.gradlew)

	return &exec.Cmd{
		Path: path,
		Args: []string{cf.gradlew, "--init-script", cf.initScript, "debrickedAllDeps"},
		Dir:  workingDirectory,
	}, err
}
