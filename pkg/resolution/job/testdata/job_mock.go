package testdata

import "fmt"

type JobMock struct {
	file string
	err  error
}

func (j *JobMock) File() string {
	return j.file
}

func (j *JobMock) Error() error {
	return j.err
}

func (j *JobMock) Run() {
	fmt.Println("job mock run")
}

func NewJobMock(file string) *JobMock {
	return &JobMock{file: file}
}

func (j *JobMock) SetErr(err error) {
	j.err = err
}
