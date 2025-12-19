package sbt

import (
	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/maven"
)

type Strategy struct {
	files           []string
	cmdFactory      ICmdFactory
	buildService    IBuildService
	mavenPomService maven.IPomService
	mavenCmdFactory maven.ICmdFactory
}

func NewStrategy(files []string) Strategy {
	return Strategy{
		files:           files,
		cmdFactory:      CmdFactory{},
		buildService:    BuildService{},
		mavenPomService: maven.PomService{},
		mavenCmdFactory: maven.CmdFactory{},
	}
}

func (s Strategy) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob

	for _, file := range s.files {
		jobs = append(jobs, NewJob(file, s.cmdFactory, s.buildService, s.mavenPomService, s.mavenCmdFactory))
	}

	return jobs, nil
}
