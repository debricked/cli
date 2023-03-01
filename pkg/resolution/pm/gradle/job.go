package gradle

import (
	"path/filepath"

	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/pm/util"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

const (
	fileName = ".debricked-gradle-dependencies.txt"
)

type Job struct {
	job.BaseJob
	cmdFactory ICmdFactory
	fileWriter writer.IFileWriter
}

func NewJob(
	file string,
	cmdFactory ICmdFactory,
	fileWriter writer.IFileWriter,
) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(file),
		cmdFactory: cmdFactory,
		fileWriter: fileWriter,
	}
}

func (j *Job) Run() {
	workingDirectory := filepath.Dir(filepath.Clean(j.GetFile()))

	dependenciesCmd, err := j.cmdFactory.MakeDependenciesCmd(workingDirectory)
	if err != nil {
		j.Errors().Critical(err)

		return
	}

	j.SendStatus("creating dependency graph")
	output, err := dependenciesCmd.Output()
	if err != nil {
		j.Errors().Critical(err)

		return
	}

	j.SendStatus("creating lock file")
	lockFile, err := j.fileWriter.Create(util.MakePathFromManifestFile(j.GetFile(), fileName))
	if err != nil {
		j.Errors().Critical(err)

		return
	}
	defer util.CloseFile(j, j.fileWriter, lockFile)

	err = j.fileWriter.Write(lockFile, output)
	if err != nil {
		j.Errors().Critical(err)
	}
}
