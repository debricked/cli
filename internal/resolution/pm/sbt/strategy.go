package sbt

import (
	"github.com/debricked/cli/internal/resolution/job"
)

type Strategy struct {
	files        []string
	cmdFactory   ICmdFactory
	buildService IBuildService
}

func NewStrategy(files []string) Strategy {
	return Strategy{
		files:        files,
		cmdFactory:   CmdFactory{},
		buildService: BuildService{},
	}
}

func (s Strategy) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob

	for _, file := range s.files {
		jobs = append(jobs, NewJob(file, s.cmdFactory, s.buildService))
	}

	return jobs, nil
}
