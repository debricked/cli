package job

import (
	"errors"
	"os/exec"
	"testing"

	err "github.com/debricked/cli/internal/io/err"
	"github.com/stretchr/testify/assert"
)

const testDir = "dir"

var testFiles = []string{"file"}

func TestNewBaseJob(t *testing.T) {
	j := NewBaseJob(testDir, testFiles)
	assert.Equal(t, testFiles, j.GetFiles())
	assert.Equal(t, testDir, j.GetDir())
	assert.NotNil(t, j.Errors())
	assert.NotNil(t, j.status)
}

func TestGetFiles(t *testing.T) {
	j := BaseJob{}
	j.files = testFiles
	assert.Equal(t, testFiles, j.GetFiles())
}

func TestGetDir(t *testing.T) {
	j := BaseJob{}
	j.dir = testDir
	assert.Equal(t, testDir, j.GetDir())
}

func TestReceiveStatus(t *testing.T) {
	j := BaseJob{
		files:  testFiles,
		dir:    testDir,
		errs:   nil,
		status: make(chan string),
	}

	statusChan := j.ReceiveStatus()
	assert.NotNil(t, statusChan)
}

func TestErrors(t *testing.T) {
	jobErr := errors.New("error")
	j := BaseJob{}
	j.dir = testDir
	j.errs = err.NewErrors(j.dir)
	j.errs.Critical(jobErr)

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), jobErr)
}

func TestSendStatus(t *testing.T) {
	j := BaseJob{
		files:  testFiles,
		dir:    testDir,
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
	differentDir := "testDifferentDir"
	differentFiles := []string{"testDifferentFile"}
	j := NewBaseJob(differentDir, differentFiles)
	assert.NotEqual(t, testFiles, j.GetFiles())
	assert.Equal(t, differentFiles, j.GetFiles())
	assert.NotEqual(t, testDir, j.GetDir())
	assert.Equal(t, differentDir, j.GetDir())
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
