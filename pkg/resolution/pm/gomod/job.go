package gomod

import (
	"os"

	"github.com/debricked/cli/pkg/resolution/pm/util"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

const (
	fileName = ".debricked-go-dependencies.txt"
)

type Job struct {
	file       string
	cmdFactory ICmdFactory
	fileWriter writer.IFileWriter
	err        error
	status     chan string
}

func NewJob(
	file string,
	cmdFactory ICmdFactory,
	fileWriter writer.IFileWriter,
) *Job {
	return &Job{
		file:       file,
		cmdFactory: cmdFactory,
		fileWriter: fileWriter,
		status:     make(chan string),
	}
}

func (j *Job) File() string {
	return j.file
}

func (j *Job) Error() error {
	return j.err
}

func (j *Job) Status() chan string {
	return j.status
}

func (j *Job) Run() {
	j.status <- "tidy dependency graph"
	_, err := j.runTidyCmd()

	// TODO Set when failing tidy to be a warning!
	if err != nil {
		j.err = err
		return
	}

	j.status <- "creating dependency graph"
	graphCmdOutput, err := j.runGraphCmd()
	if err != nil {
		j.err = err
		return
	}

	j.status <- "creating dependency version list"
	listCmdOutput, err := j.runListCmd()
	if err != nil {
		j.err = err
		return
	}

	j.status <- "creating lock file"
	lockFile, err := j.fileWriter.Create(util.MakePathFromManifestFile(j.file, fileName))
	if err != nil {
		j.err = err

		return
	}
	defer closeFile(j, lockFile)

	var fileContents []byte
	fileContents = append(fileContents, graphCmdOutput...)
	fileContents = append(fileContents, []byte("\n")...)
	fileContents = append(fileContents, listCmdOutput...)

	j.err = j.fileWriter.Write(lockFile, fileContents)
}

func (j *Job) runTidyCmd() ([]byte, error) {
	tidyCmd, err := j.cmdFactory.MakeTidyCmd()
	if err != nil {
		j.err = err

		return nil, err
	}

	tidyCmdOutput, err := tidyCmd.Output()
	if err != nil {
		j.err = err

		return nil, err
	}

	return tidyCmdOutput, nil
}

func (j *Job) runGraphCmd() ([]byte, error) {
	graphCmd, err := j.cmdFactory.MakeGraphCmd()
	if err != nil {
		j.err = err

		return nil, err
	}

	graphCmdOutput, err := graphCmd.Output()
	if err != nil {
		j.err = err

		return nil, err
	}

	return graphCmdOutput, nil
}

func (j *Job) runListCmd() ([]byte, error) {
	listCmd, err := j.cmdFactory.MakeListCmd()
	if err != nil {
		j.err = err

		return nil, err
	}

	listCmdOutput, err := listCmd.Output()
	if err != nil {
		j.err = err

		return nil, err
	}

	return listCmdOutput, nil
}

func closeFile(job *Job, file *os.File) {
	err := job.fileWriter.Close(file)
	if err != nil {
		job.err = err
	}
}
