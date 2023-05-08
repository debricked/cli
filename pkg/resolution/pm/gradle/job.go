package gradle

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

type Job struct {
	job.BaseJob
	dir              string
	gradlew          string
	groovyInitScript string
	cmdFactory       ICmdFactory
	fileWriter       writer.IFileWriter
}

func NewJob(
	file string,
	dir string,
	gradlew string,
	groovyInitScript string,
	cmdFactory ICmdFactory,
	fileWriter writer.IFileWriter,
) *Job {

	return &Job{
		BaseJob:          job.NewBaseJob(file),
		dir:              dir,
		gradlew:          gradlew,
		groovyInitScript: groovyInitScript,
		cmdFactory:       cmdFactory,
		fileWriter:       fileWriter,
	}
}

func (j *Job) Run() {
	workingDirectory := filepath.Clean(j.GetDir())
	dependenciesCmd, err := j.cmdFactory.MakeDependenciesGraphCmd(workingDirectory, j.gradlew, j.groovyInitScript)
	var permissionErr error

	if err != nil {
		if strings.HasSuffix(err.Error(), "gradlew\": permission denied") {
			permissionDeniedErr := "Permission to execute gradlew is not granted, fallback to PATHs gradle installation will be used.\nFull error: " + err.Error()
			dependenciesCmd, err = j.cmdFactory.MakeDependenciesGraphCmd(workingDirectory, "gradle", j.groovyInitScript)
			permissionErr = errors.New(permissionDeniedErr)
		}
	}

	if err != nil {
		if permissionErr != nil {
			j.Errors().Critical(permissionErr)
		}
		j.Errors().Critical(err)
		return
	}

	j.SendStatus("creating dependency graph")
	_, err = dependenciesCmd.Output()

	if permissionErr != nil {
		if err != nil {
			j.Errors().Critical(permissionErr)
		} else {
			j.Errors().Warning(permissionErr)
		}
	}

	if err != nil {
		j.Errors().Critical(j.GetExitError(err))

		return
	}
}

func (j *Job) GetDir() string {
	return j.dir
}
