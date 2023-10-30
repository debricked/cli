package testdata

import (
	"fmt"

	"github.com/debricked/cli/internal/io/err"
)

type JobMock struct {
	dir    string
	files  []string
	errs   err.IErrors
	status chan string
}

func (j *JobMock) ReceiveStatus() chan string {
	return j.status
}

func (j *JobMock) GetDir() string {
	return j.dir
}

func (j *JobMock) GetFiles() []string {
	return j.files
}

func (j *JobMock) Errors() err.IErrors {
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
		errs:   err.NewErrors(dir),
	}
}

func (j *JobMock) SetErr(err err.IError) {
	j.errs.Critical(err)
}
