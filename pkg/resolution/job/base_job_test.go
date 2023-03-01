package job

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFile(t *testing.T) {
	j := BaseJob{}
	j.File = "file"
	assert.Equal(t, "file", j.GetFile())
}

func TestReceiveStatus(t *testing.T) {
	j := BaseJob{
		File:   "file",
		Err:    nil,
		Status: make(chan string),
	}

	statusChan := j.ReceiveStatus()
	assert.NotNil(t, statusChan)
}

func TestError(t *testing.T) {
	jobErr := errors.New("error")
	j := BaseJob{}
	j.File = "file"
	j.Err = jobErr
	assert.Equal(t, jobErr, j.Error())
}

func TestSendStatus(t *testing.T) {
	j := BaseJob{
		File:   "file",
		Err:    nil,
		Status: make(chan string),
	}

	go func() {
		status := <-j.ReceiveStatus()
		assert.Equal(t, "status", status)
	}()

	j.SendStatus("status")
}
