package job

import (
	"errors"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testFile = "file"

func TestNewBaseJob(t *testing.T) {
	j := NewBaseJob(testFile)
	assert.Equal(t, testFile, j.GetFile())
	assert.NotNil(t, j.Errors())
	assert.NotNil(t, j.status)
}

func TestGetFile(t *testing.T) {
	j := BaseJob{}
	j.file = testFile
	assert.Equal(t, testFile, j.GetFile())
}

func TestReceiveStatus(t *testing.T) {
	j := BaseJob{
		file:   testFile,
		errs:   nil,
		status: make(chan string),
	}

	statusChan := j.ReceiveStatus()
	assert.NotNil(t, statusChan)
}

func TestErrors(t *testing.T) {
	jobErr := errors.New("error")
	j := BaseJob{}
	j.file = testFile
	j.errs = NewErrors(j.file)
	jobError := NewBaseJobError(jobErr.Error())
	j.errs.Critical(jobError)

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), jobError)
}

func TestSendStatus(t *testing.T) {
	j := BaseJob{
		file:   testFile,
		errs:   nil,
		status: make(chan string),
	}

	go func() {
		status := <-j.ReceiveStatus()
		assert.Equal(t, "status", status)
	}()

	j.SendStatus("status")
}

func TestDifferentNewBaseJob(t *testing.T) {
	differentFileName := "testDifferentFile"
	j := NewBaseJob(differentFileName)
	assert.NotEqual(t, testFile, j.GetFile())
	assert.Equal(t, differentFileName, j.GetFile())
	assert.NotNil(t, j.Errors())
	assert.NotNil(t, j.status)
}

func TestGetExitErrorWithExitError(t *testing.T) {
	err := &exec.ExitError{
		ProcessState: nil,
		Stderr:       []byte("stderr"),
	}
	j := BaseJob{}
	exitErr := j.GetExitError(err)
	assert.ErrorContains(t, exitErr, string(err.Stderr))
}

func TestGetExitErrorWithNoneExitError(t *testing.T) {
	err := &exec.Error{Err: errors.New("none-exit-err")}
	j := BaseJob{}
	exitErr := j.GetExitError(err)
	assert.ErrorContains(t, exitErr, err.Error())
}
