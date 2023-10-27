package maven

import (
	"errors"
	"path/filepath"

	"github.com/debricked/cli/internal/resolution/job"
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
		j.Errors().Critical(err)

		return
	}
	j.SendStatus("creating dependency graph")
	var output []byte
	output, err = cmd.Output()
	if err != nil {
		if output == nil {
			j.Errors().Critical(err)
		} else {
			j.Errors().Critical(errors.New(string(output)))
		}
	}
}
