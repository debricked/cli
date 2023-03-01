package gradle

import (
	"path/filepath"

	"github.com/debricked/cli/pkg/resolution/job"
	jobTestdata "github.com/debricked/cli/pkg/resolution/job/testdata"
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
		BaseJob: job.BaseJob{
			File:   file,
			Status: make(chan string),
		},
		cmdFactory: cmdFactory,
		fileWriter: fileWriter,
	}
}

func (j *Job) Run() {
	workingDirectory := filepath.Dir(filepath.Clean(j.GetFile()))

	dependenciesCmd, err := j.cmdFactory.MakeDependenciesCmd(workingDirectory)
	if err != nil {
		j.Err = err

		return
	}

	j.SendStatus("creating dependency graph")
	output, err := dependenciesCmd.Output()
	if err != nil {
		j.Err = err

		return
	}

	j.SendStatus("creating lock file")
	lockFile, err := j.fileWriter.Create(util.MakePathFromManifestFile(j.GetFile(), fileName))
	if err != nil {
		j.Err = err

		return
	}
	defer jobTestdata.CloseFile(&j.BaseJob, j.fileWriter, lockFile)

	j.Err = j.fileWriter.Write(lockFile, output)
}
