package maven

import "os/exec"

type ICmdFactory interface {
	MakeDependencyTreeCmd() (*exec.Cmd, error)
}

type CmdFactory struct{}

func (_ CmdFactory) MakeDependencyTreeCmd() (*exec.Cmd, error) {
	path, err := exec.LookPath("mvn")

	return &exec.Cmd{
		Path: path,
		Args: []string{"mvn", "dependency:tree", "-DoutputFile=.debricked-maven-dependencies.tgf", "-DoutputType=tgf"},
	}, err
}
