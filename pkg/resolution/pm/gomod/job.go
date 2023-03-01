package gomod

import (
	"github.com/debricked/cli/pkg/resolution/job"
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
		BaseJob:    job.NewBaseJob(file),
		cmdFactory: cmdFactory,
		fileWriter: fileWriter,
	}
}

func (j *Job) Run() {
	j.SendStatus("creating dependency graph")
	graphCmdOutput, err := j.runGraphCmd()
	if err != nil {
		j.Errors().Critical(err)

		return
	}

	j.SendStatus("creating dependency version list")
	listCmdOutput, err := j.runListCmd()
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

	var fileContents []byte
	fileContents = append(fileContents, graphCmdOutput...)
	fileContents = append(fileContents, []byte("\n")...)
	fileContents = append(fileContents, listCmdOutput...)

	err = j.fileWriter.Write(lockFile, fileContents)
	if err != nil {
		j.Errors().Critical(err)
	}
}

func (j *Job) runGraphCmd() ([]byte, error) {
	graphCmd, err := j.cmdFactory.MakeGraphCmd()
	if err != nil {
		return nil, err
	}

	graphCmdOutput, err := graphCmd.Output()
	if err != nil {
		return nil, err
	}

	return graphCmdOutput, nil
}

func (j *Job) runListCmd() ([]byte, error) {
	listCmd, err := j.cmdFactory.MakeListCmd()
	if err != nil {
		return nil, err
	}

	listCmdOutput, err := listCmd.Output()
	if err != nil {
		return nil, err
	}

	return listCmdOutput, nil
}
