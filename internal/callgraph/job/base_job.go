package job

import (
	"errors"
	"os/exec"

	err "github.com/debricked/cli/internal/io/err"
)

type BaseJob struct {
	dir    string
	files  []string
	errs   err.IErrors
	status chan string
}

func NewBaseJob(dir string, files []string) BaseJob {
	return BaseJob{
		dir:    dir,
		files:  files,
		errs:   err.NewErrors(dir),
		status: make(chan string),
	}
}

func (j *BaseJob) GetDir() string {
	return j.dir
}

func (j *BaseJob) GetFiles() []string {
	return j.files
}

func (j *BaseJob) Errors() err.IErrors {
	return j.errs
}

func (j *BaseJob) ReceiveStatus() chan string {
	return j.status
}

func (j *BaseJob) SendStatus(status string) {
	j.status <- status
}

func (j *BaseJob) GetExitError(err error) error {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return err
	}

	return errors.New(string(exitErr.Stderr))
}
