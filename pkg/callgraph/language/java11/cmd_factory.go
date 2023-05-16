package java

import (
	"os/exec"
)

type ICmdFactory interface {
	MakeMvnCopyDependenciesCmd(workingDirectory string, targetDir string) (*exec.Cmd, error)
	MakeCallGraphGenerationCmd(callgraphJarPath string, workingDirectory string, targetClasses string, dependencyClasses string) (*exec.Cmd, error)
}

type CmdFactory struct{}

func (_ CmdFactory) MakeMvnCopyDependenciesCmd(
	workingDirectory string,
	targetDir string,
) (*exec.Cmd, error) {
	path, err := exec.LookPath("mvn")

	args := []string{
		"mvn",
		"-q",
		"-B",
		"dependency:copy-dependencies",
		"-DoutputDirectory=" + targetDir,
		"-DskipTests",
	}

	return &exec.Cmd{
		Path: path,
		Args: args,
		Dir:  workingDirectory,
	}, err
}

func (_ CmdFactory) MakeCallGraphGenerationCmd(
	callgraphJarPath string,
	workingDirectory string,
	targetClasses string,
	dependencyClasses string,
) (*exec.Cmd, error) {
	path, err := exec.LookPath("java")
	args := []string{
		"java",
		"-jar",
		callgraphJarPath,
		"-u",
		targetClasses,
		"-l",
		dependencyClasses,
		"-f",
		".debricked-call-graph",
	}

	return &exec.Cmd{
		Path: path,
		Args: args,
		Dir:  workingDirectory,
	}, err
}
