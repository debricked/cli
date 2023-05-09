package java

import "os/exec"

type ICmdFactory interface {
	MakeMvnCopyDependenciesCmd(workingDirectory string, targetRootPomDir string) (*exec.Cmd, error)
	MakeGradleCopyDependenciesCmd(workingDirectory string, targetRootPomDir string) (*exec.Cmd, error)
	MakeCallGraphGenerationCmd(callgraphJarPath string, workingDirectory string, targetClasses string, dependencyClasses string) (*exec.Cmd, error)
}

type CmdFactory struct{}

func (_ CmdFactory) MakeGradleCopyDependenciesCmd(
	workingDirectory string,
	gradlew string,
) (*exec.Cmd, error) {
	path, err := exec.LookPath(gradlew)

	// TargetDir already in groovy script
	return &exec.Cmd{
		Path: path,
		Args: []string{
			"gradle",
			"-q",
			"debrickedCopyDependencies",
		},
		Dir: workingDirectory,
	}, err
}

func (_ CmdFactory) MakeMvnCopyDependenciesCmd(
	workingDirectory string,
	targetDir string,
) (*exec.Cmd, error) {
	path, err := exec.LookPath("mvn")

	return &exec.Cmd{
		Path: path,
		Args: []string{
			"mvn",
			"-q",
			"-B",
			"dependency:copy-dependencies",
			"-DoutputDirectory=" + targetDir,
			"-DskipTests",
		},
		Dir: workingDirectory,
	}, err
}

func (_ CmdFactory) MakeCallGraphGenerationCmd(
	callgraphJarPath string,
	workingDirectory string,
	targetClasses string,
	dependencyClasses string,
) (*exec.Cmd, error) {
	path, err := exec.LookPath("java")

	return &exec.Cmd{
		Path: path,
		Args: []string{
			"java",
			"-jar",
			callgraphJarPath,
			"-u",
			targetClasses,
			"-l",
			dependencyClasses,
			"-f",
			".debricked-call-graph",
		},
		Dir: workingDirectory,
	}, err
}
