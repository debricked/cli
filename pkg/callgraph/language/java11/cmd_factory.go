package java

import "os/exec"

type ICmdFactory interface {
	MakeCallGraphGenerationCmd(workingDirectory string, targetClasses string, dependencyClasses string) (*exec.Cmd, error)
}

type CmdFactory struct{}

func (_ CmdFactory) MakeCallGraphGenerationCmd(
	workingDirectory string,
	targetClasses string,
	dependencyClasses string,
) (*exec.Cmd, error) {
	path, err := exec.LookPath("java")
	javaVulnFunc := "/home/magnus/Projects/debricked/vulnerable-functionality-github/vulnerable-functionality/java/common/target/SootWrapper.jar"

	return &exec.Cmd{
		Path: path,
		Args: []string{
			"java",
			"-jar",
			javaVulnFunc,
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
