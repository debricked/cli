package nuget

import (
	"errors"

	"github.com/debricked/cli/internal/resolution/job"
)

const (
	nuget = "dotnet"
)

type Job struct {
	job.BaseJob
	install      bool
	nugetCommand string
	cmdFactory   ICmdFactory
}

func NewJob(
	file string,
	install bool,
	cmdFactory ICmdFactory,
) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(file),
		install:    install,
		cmdFactory: cmdFactory,
	}
}

func (j *Job) Install() bool {
	return j.install
}

func (j *Job) Run() {
	if j.install {

		j.SendStatus("installing dependencies")
		output, err := j.runInstallCmd()
		if err != nil {
			if output == nil {
				j.Errors().Critical(err)
			} else {
				j.Errors().Critical(errors.New(string(output)))
			}

			return
		}
	}

}

func (j *Job) runInstallCmd() ([]byte, error) {

	j.nugetCommand = nuget
	installCmd, err := j.cmdFactory.MakeInstallCmd(j.nugetCommand, j.GetFile())
	if err != nil {
		return nil, err
	}

	installCmdOutput, err := installCmd.Output()
	if err != nil {
		return installCmdOutput, j.GetExitError(err)
	}

	return installCmdOutput, nil
}
