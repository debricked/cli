package pip

import (
	"errors"
	"testing"

	"github.com/debricked/cli/pkg/resolution/pm/writer"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	job := NewJob("file", false, CmdFactory{}, writer.FileWriter{})
	assert.Equal(t, "file", job.file)
	assert.Nil(t, job.err)
}

func TestFile(t *testing.T) {
	job := Job{file: "file"}
	assert.Equal(t, "file", job.File())
}

func TestError(t *testing.T) {
	jobErr := errors.New("error")
	job := Job{file: "file", err: jobErr}
	assert.Equal(t, jobErr, job.Error())
}

// TODO add more tests a la maven / golang
