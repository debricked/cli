package job

import (
	"errors"
	"os/exec"
	"strings"
)

type BaseJob struct {
	file   string
	errs   IErrors
	status chan string
}

func NewBaseJob(file string) BaseJob {
	return BaseJob{
		file:   file,
		errs:   NewErrors(file),
		status: make(chan string),
	}
}

func (j *BaseJob) GetFile() string {
	return j.file
}

func (j *BaseJob) Errors() IErrors {
	return j.errs
}

func (j *BaseJob) ReceiveStatus() chan string {
	return j.status
}

func (j *BaseJob) SendStatus(status string) {
	j.status <- status
}

func (j *BaseJob) GetExitError(err error, commandOutput string) error {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return err
	}

	// If Stderr is empty, use commandOutput as error string instead
	errorMessage := string(exitErr.Stderr)
	if errorMessage == "" {
		errorMessage = commandOutput
	}

	return errors.New(errorMessage)
}

func (j *BaseJob) GetExecutableNotFoundErrorDocumentation(pm string) string {
	return strings.Join(
		[]string{
			pm + " wasn't found.",
			"Please check if it is installed and accessible by the CLI.",
		}, " ")
}
