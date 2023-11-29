package composer

import (
	"github.com/debricked/cli/internal/resolution/job"
)

const (
	composer = "composer"
)

type Job struct {
	job.BaseJob
	install         bool
	composerCommand string
	cmdFactory      ICmdFactory
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
		_, err := j.runInstallCmd()
		if err != nil {
			j.Errors().Critical(err)

			return
		}
	}

}

func (j *Job) runInstallCmd() ([]byte, error) {

	j.composerCommand = composer
	installCmd, err := j.cmdFactory.MakeInstallCmd(j.composerCommand, j.GetFile())
	if err != nil {
		return nil, err
	}

	installCmdOutput, err := installCmd.Output()
	if err != nil {
		return nil, j.GetExitError(err)
	}

	return installCmdOutput, nil
}
