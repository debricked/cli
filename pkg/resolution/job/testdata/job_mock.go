package testdata

import (
	"fmt"

	"github.com/debricked/cli/pkg/resolution/job"
)

type JobMock struct {
	file   string
	errs   job.IErrors
	status chan string
}

func (j *JobMock) ReceiveStatus() chan string {
	return j.status
}

func (j *JobMock) GetFile() string {
	return j.file
}

func (j *JobMock) Errors() job.IErrors {
	return j.errs
}

func (j *JobMock) Run() {
	fmt.Println("job mock run")
}

func NewJobMock(file string) *JobMock {
	return &JobMock{
		file:   file,
		status: make(chan string),
		errs:   job.NewErrors(file),
	}
}

func (j *JobMock) SetErr(err job.IError) {
	j.errs.Critical(err)
}
