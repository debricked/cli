package gradle

import (
	"errors"
	"github.com/debricked/cli/pkg/resolution/pm/gradle/testdata"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
	writerTestdata "github.com/debricked/cli/pkg/resolution/pm/writer/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJob(t *testing.T) {
	job := NewJob("file", CmdFactory{}, writer.FileWriter{})
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
	job := NewJob("file", testdata.CmdFactoryMock{Err: cmdErr}, nil)
	job.Run()

	assert.ErrorIs(T, cmdErr, job.Error())
}

func TestRunCmdOutputErr(T *testing.T) {
	job := NewJob("file", testdata.CmdFactoryMock{Name: "bad-name"}, nil)
	job.Run()

	assert.ErrorContains(T, job.err, "executable file not found in $PATH")
}

func TestRunCreateErr(T *testing.T) {
	createErr := errors.New("create-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}
	job := NewJob("file", testdata.CmdFactoryMock{Name: "echo"}, fileWriterMock)
	job.Run()

	assert.ErrorIs(T, job.Error(), createErr)
}

func TestRunWriteErr(T *testing.T) {
	writeErr := errors.New("write-error")
	fileWriterMock := &writerTestdata.FileWriterMock{WriteErr: writeErr}
	job := NewJob("file", testdata.CmdFactoryMock{Name: "echo"}, fileWriterMock)
	job.Run()

	assert.ErrorIs(T, job.Error(), writeErr)
}

func TestRunCloseErr(T *testing.T) {
	closeErr := errors.New("close-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CloseErr: closeErr}
	job := NewJob("file", testdata.CmdFactoryMock{Name: "echo"}, fileWriterMock)
	job.Run()

	assert.ErrorIs(T, job.Error(), closeErr)
}

func TestRun(T *testing.T) {
	fileContents := []byte("MakeDependenciesCmd\n")
	fileWriterMock := &writerTestdata.FileWriterMock{}
	cmdFactoryMock := testdata.CmdFactoryMock{Name: "echo"}
	job := NewJob("file", cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.NoError(T, job.Error())
	assert.Equal(T, fileContents, fileWriterMock.Contents)
}
