package gomod

import (
	"path/filepath"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/debricked/cli/internal/resolution/pm/writer"
)

const (
	fileName = "gomod.debricked.lock"
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

	workingDirectory := filepath.Dir(filepath.Clean(j.GetFile()))

	graphCmdOutput, err := j.runGraphCmd(workingDirectory)
	if err != nil {
		j.Errors().Critical(util.NewPMJobError(err.Error()))

		return
	}

	j.SendStatus("creating dependency version list")
	listCmdOutput, err := j.runListCmd(workingDirectory)
	if err != nil {
		j.Errors().Critical(util.NewPMJobError(err.Error()))

		return
	}

	j.SendStatus("creating lock file")
	lockFile, err := j.fileWriter.Create(util.MakePathFromManifestFile(j.GetFile(), fileName))
	if err != nil {
		j.Errors().Critical(util.NewPMJobError(err.Error()))

		return
	}
	defer util.CloseFile(j, j.fileWriter, lockFile)

	var fileContents []byte
	fileContents = append(fileContents, graphCmdOutput...)
	fileContents = append(fileContents, []byte("\n")...)
	fileContents = append(fileContents, listCmdOutput...)

	err = j.fileWriter.Write(lockFile, fileContents)
	if err != nil {
		j.Errors().Critical(util.NewPMJobError(err.Error()))
	}
}

func (j *Job) runGraphCmd(workingDirectory string) ([]byte, error) {
	graphCmd, err := j.cmdFactory.MakeGraphCmd(workingDirectory)
	if err != nil {
		return nil, err
	}

	graphCmdOutput, err := graphCmd.Output()
	if err != nil {
		return nil, j.GetExitError(err)
	}

	return graphCmdOutput, nil
}

func (j *Job) runListCmd(workingDirectory string) ([]byte, error) {
	listCmd, err := j.cmdFactory.MakeListCmd(workingDirectory)
	if err != nil {
		return nil, err
	}

	listCmdOutput, err := listCmd.Output()
	if err != nil {
		return nil, j.GetExitError(err)
	}

	return listCmdOutput, nil
}
