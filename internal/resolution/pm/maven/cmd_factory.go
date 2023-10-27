package maven

import "os/exec"

type ICmdFactory interface {
	MakeDependencyTreeCmd(workingDirectory string) (*exec.Cmd, error)
}

type CmdFactory struct{}

func (_ CmdFactory) MakeDependencyTreeCmd(workingDirectory string) (*exec.Cmd, error) {
	path, err := exec.LookPath("mvn")

	return &exec.Cmd{
		Path: path,
		Args: []string{
			"mvn",
			"dependency:tree",
			"-DoutputFile=" + lockFileExtension,
			"-DoutputType=tgf",
			"--fail-at-end",
		},
		Dir: workingDirectory,
	}, err
}
