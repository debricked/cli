package testdata

import (
	"fmt"

	"github.com/debricked/cli/pkg/resolution/job"
)

type JobMock struct {
	file   string
	err    error
	status chan string
}

func (j *JobMock) ReceiveStatus() chan string {
	return j.status
}

func (j *JobMock) GetFile() string {
	return j.file
}

func (j *JobMock) Error() job.IJobError {
	return j.err
}

func (j *JobMock) Run() {
	fmt.Println("job mock run")
}

func NewJobMock(file string) *JobMock {
	return &JobMock{file: file, status: make(chan string)}
}

func (j *JobMock) SetErr(err error) {
	j.err = err
}
