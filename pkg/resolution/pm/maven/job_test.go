package maven

import (
	"errors"
	"testing"

	"github.com/debricked/cli/pkg/resolution/pm/maven/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	job := NewJob("file", CmdFactory{})
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

func TestRunCmdErr(T *testing.T) {
	cmdErr := errors.New("cmd-error")
	job := NewJob("file", testdata.CmdFactoryMock{Err: cmdErr})

	job.Run()

	assert.ErrorIs(T, cmdErr, job.Error())
}

func TestRunCmdOutputErr(T *testing.T) {
	job := NewJob("file", testdata.CmdFactoryMock{Name: "bad-name"})

	job.Run()

	assert.ErrorContains(T, job.err, "executable file not found in")
	assert.ErrorContains(T, job.err, "PATH")
}

func TestRun(T *testing.T) {
	job := NewJob("file", testdata.CmdFactoryMock{Name: "echo"})

	job.Run()

	assert.NoError(T, job.err)
}
