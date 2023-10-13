package job

import (
	"errors"
	"fmt"
	"os/exec"
)

type BaseJob struct {
	file          string
	errs          IErrors
	status        chan string
	currentStatus string
}

func NewBaseJob(file string) BaseJob {
	return BaseJob{
		file:          file,
		errs:          NewErrors(file),
		status:        make(chan string),
		currentStatus: "",
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

func (j *BaseJob) FmtError(err error, output []byte) error {

	if err == nil {
		return fmt.Errorf("%s error: No error was present.", j.currentStatus)
	}

	errorString := fmt.Errorf("%s error: %s", j.currentStatus, err)

	if output != nil {
		errorString = fmt.Errorf("%s output: %s\n%s", j.currentStatus, output, errorString)
	}

	return errorString
}

func (j *BaseJob) SendStatus(status string) {
	j.currentStatus = status
	j.status <- status
}

func (j *BaseJob) GetExitError(err error) error {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return err
	}

	return errors.New(string(exitErr.Stderr))
}
