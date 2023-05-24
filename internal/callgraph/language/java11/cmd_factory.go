package java

import (
	"os/exec"

	"github.com/debricked/cli/internal/callgraph/cgexec"
)

type ICmdFactory interface {
	MakeMvnCopyDependenciesCmd(workingDirectory string, targetDir string, ctx cgexec.IContext) (*exec.Cmd, error)
	MakeCallGraphGenerationCmd(callgraphJarPath string, workingDirectory string, targetClasses string, dependencyClasses string, ctx cgexec.IContext) (*exec.Cmd, error)
	MakeBuildMavenCmd(workingDirectory string, ctx cgexec.IContext) (*exec.Cmd, error)
}

type CmdFactory struct{}

func (_ CmdFactory) MakeMvnCopyDependenciesCmd(
	workingDirectory string,
	targetDir string,
	ctx cgexec.IContext,
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

	return cgexec.MakeCommand(workingDirectory, path, args, ctx), err
}

func (_ CmdFactory) MakeCallGraphGenerationCmd(
	callgraphJarPath string,
	workingDirectory string,
	targetClasses string,
	dependencyClasses string,
	ctx cgexec.IContext,
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

	return cgexec.MakeCommand(workingDirectory, path, args, ctx), err
}

func (_ CmdFactory) MakeBuildMavenCmd(workingDirectory string, ctx cgexec.IContext) (*exec.Cmd, error) {
	// NOTE: mvn compile should work in theory and be faster
	path, err := exec.LookPath("mvn")
	args := []string{
		"mvn",
		"package",
		"-q",
		"-DskipTests",
	}
	return cgexec.MakeCommand(workingDirectory, path, args, ctx), err
}
