package maven

import (
	"path/filepath"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	lockFileExtension = "maven.debricked.lock"
)

type Job struct {
	job.BaseJob
	cmdFactory ICmdFactory
}

func NewJob(file string, cmdFactory ICmdFactory) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(file),
		cmdFactory: cmdFactory,
	}
}

func (j *Job) Run() {
	workingDirectory := filepath.Dir(filepath.Clean(j.GetFile()))
	cmd, err := j.cmdFactory.MakeDependencyTreeCmd(workingDirectory)
	if err != nil {
		j.Errors().Critical(util.NewPMJobError(err.Error()))

		return
	}
	j.SendStatus("creating dependency graph")
	var output []byte
	output, err = cmd.Output()
	if err != nil {
		if output == nil {
			j.Errors().Critical(util.NewPMJobError(err.Error()))
		} else {
			j.Errors().Critical(util.NewPMJobError(string(output)))
		}
	}
}
