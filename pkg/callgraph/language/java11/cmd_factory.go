package java

import "os/exec"

type ICmdFactory interface {
	MakeBuildMvnCopyDependenciesCmd(workingDirectory string, targetRootPomDir string) (*exec.Cmd, error)
	MakeBuildGradleCopyDependenciesCmd(workingDirectory string, targetRootPomDir string) (*exec.Cmd, error)
	MakeCallGraphGenerationCmd(callgraphJarPath string, workingDirectory string, targetClasses string, dependencyClasses string) (*exec.Cmd, error)
}

type CmdFactory struct{}

func (_ CmdFactory) MakeBuildGradleCopyDependenciesCmd(
	workingDirectory string,
	targetRootPomDir string,
) (*exec.Cmd, error) {
	path, err := exec.LookPath("mvn")

	return &exec.Cmd{
		Path: path,
		Args: []string{
			"mvn",
			"-q",
			"-B",
			"-f",
			"pom.xml",
			"package",
			"dependency:copy-dependencies",
			"-DoutputDirectory=" + targetRootPomDir,
			"-DskipTests",
		},
		Dir: workingDirectory,
	}, err
}

func (_ CmdFactory) MakeBuildMvnCopyDependenciesCmd(
	workingDirectory string,
	targetRootPomDir string,
) (*exec.Cmd, error) {
	path, err := exec.LookPath("mvn")

	return &exec.Cmd{
		Path: path,
		Args: []string{
			"mvn",
			"-q",
			"-B",
			"-f",
			"pom.xml",
			"package",
			"dependency:copy-dependencies",
			"-DoutputDirectory=" + targetRootPomDir,
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
