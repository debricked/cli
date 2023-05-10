package testdata

import (
	"fmt"

	"github.com/debricked/cli/pkg/resolution/job"
)

type JobMock struct {
	dir    string
	files  []string
	errs   job.IErrors
	status chan string
}

func (j *JobMock) ReceiveStatus() chan string {
	return j.status
}

func (j *JobMock) GetDir() string {
	return j.dir
}

func (j *JobMock) GetFiles() string {
	return j.file
}

func (j *JobMock) Errors() job.IErrors {
	return j.errs
}

func (j *JobMock) Run() {
	fmt.Println("job mock run")
}

func NewJobMock(dir string, files []string) *JobMock {
	return &JobMock{
		dir:    dir,
		files:  files,
		status: make(chan string),
		errs:   job.NewErrors(file),
	}
}

func (j *JobMock) SetErr(err job.IError) {
	j.errs.Critical(err)
}
