package gomod

import (
	"github.com/debricked/cli/pkg/resolution/job"
	jobTestdata "github.com/debricked/cli/pkg/resolution/job/testdata"
	"github.com/debricked/cli/pkg/resolution/pm/util"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

const (
	fileName = ".debricked-go-dependencies.txt"
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
	j.SendStatus("creating dependency graph")
	graphCmdOutput, err := j.runGraphCmd()
	if err != nil {
		j.err = err

		return
	}

	j.SendStatus("creating dependency version list")
	listCmdOutput, err := j.runListCmd()
	if err != nil {
		j.err = err

		return
	}

	j.SendStatus("creating lock file")
	lockFile, err := j.fileWriter.Create(util.MakePathFromManifestFile(j.GetFile(), fileName))
	if err != nil {
		j.Err = err

		return
	}
	defer jobTestdata.CloseFile(&j.BaseJob, j.fileWriter, lockFile)

	var fileContents []byte
	fileContents = append(fileContents, graphCmdOutput...)
	fileContents = append(fileContents, []byte("\n")...)
	fileContents = append(fileContents, listCmdOutput...)

	j.Err = j.fileWriter.Write(lockFile, fileContents)
}

func (j *Job) runGraphCmd() ([]byte, error) {
	graphCmd, err := j.cmdFactory.MakeGraphCmd()
	if err != nil {
		j.Err = err

		return nil, err
	}

	graphCmdOutput, err := graphCmd.Output()
	if err != nil {
		j.Err = err

		return nil, err
	}

	return graphCmdOutput, nil
}

func (j *Job) runListCmd() ([]byte, error) {
	listCmd, err := j.cmdFactory.MakeListCmd()
	if err != nil {
		j.Err = err

		return nil, err
	}

	listCmdOutput, err := listCmd.Output()
	if err != nil {
		j.Err = err

		return nil, err
	}

	return listCmdOutput, nil
}
