package gradle

import (
	"path/filepath"

	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

const (
	fileName = ".gradle.debricked.lock"
)

type Job struct {
	job.BaseJob
	gradlew          string
	groovyInitScript string
	cmdFactory       ICmdFactory
	fileWriter       writer.IFileWriter
}

func NewJob(
	file string,
	gradlew string,
	groovyInitScript string,
	cmdFactory ICmdFactory,
	fileWriter writer.IFileWriter,
) *Job {
	return &Job{
		BaseJob:          job.NewBaseJob(file),
		gradlew:          gradlew,
		groovyInitScript: groovyInitScript,
		cmdFactory:       cmdFactory,
		fileWriter:       fileWriter,
	}
}

func (j *Job) Run() {
	workingDirectory := filepath.Clean(j.GetFile())
	dependenciesCmd, err := j.cmdFactory.MakeDependenciesGraphCmd(workingDirectory, j.gradlew, j.groovyInitScript)
	if err != nil {
		j.Errors().Critical(err)

		return
	}

	j.SendStatus("creating dependency graph")
	_, err = dependenciesCmd.Output()
	if err != nil {
		j.Errors().Critical(j.GetExitError(err))

		return
	}
}
