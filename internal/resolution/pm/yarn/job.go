package yarn

import (
	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	yarn = "yarn"
)

type Job struct {
	job.BaseJob
	install     bool
	yarnCommand string
	cmdFactory  ICmdFactory
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
			jobError := util.NewPMJobError(err.Error())
			j.Errors().Critical(jobError)

			return
		}
	}

}

func (j *Job) runInstallCmd() ([]byte, error) {

	j.yarnCommand = yarn
	installCmd, err := j.cmdFactory.MakeInstallCmd(j.yarnCommand, j.GetFile())
	if err != nil {
		return nil, err
	}

	installCmdOutput, err := installCmd.Output()
	if err != nil {
		return nil, j.GetExitError(err)
	}

	return installCmdOutput, nil
}
