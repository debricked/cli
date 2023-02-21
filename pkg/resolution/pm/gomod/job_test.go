package gomod

import (
	"errors"
	"testing"

	"github.com/debricked/cli/pkg/resolution/pm/gomod/testdata"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
	writerTestdata "github.com/debricked/cli/pkg/resolution/pm/writer/testdata"
	"github.com/stretchr/testify/assert"
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

func TestRunGraphCmdErr(T *testing.T) {
	cmdErr := errors.New("cmd-error")
	job := NewJob("file", testdata.CmdFactoryMock{MakeGraphCmdErr: cmdErr}, nil)

	job.Run()

	assert.ErrorIs(T, cmdErr, job.Error())
}

func TestRunCmdOutputErr(T *testing.T) {
	job := NewJob("file", testdata.CmdFactoryMock{GraphCmdName: "bad-name"}, nil)

	job.Run()

	assert.ErrorContains(T, job.err, "executable file not found in")
	assert.ErrorContains(T, job.err, "PATH")
}

func TestRunListCmdErr(T *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.CmdFactoryMock{
		GraphCmdName:   "echo",
		MakeListCmdErr: cmdErr,
	}
	job := NewJob("file", cmdFactoryMock, nil)

	job.Run()

	assert.ErrorIs(T, cmdErr, job.Error())
}

func TestRunListCmdOutputErr(T *testing.T) {
	job := NewJob("file", testdata.CmdFactoryMock{GraphCmdName: "echo", ListCmdName: "bad-name"}, nil)

	job.Run()

	assert.ErrorContains(T, job.err, "executable file not found in")
	assert.ErrorContains(T, job.err, "PATH")
}

func TestRunCreateErr(T *testing.T) {
	createErr := errors.New("create-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}
	cmdFactoryMock := testdata.CmdFactoryMock{GraphCmdName: "echo", ListCmdName: "echo"}
	job := NewJob("file", cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.ErrorIs(T, job.Error(), createErr)
}

func TestRunWriteErr(T *testing.T) {
	writeErr := errors.New("write-error")
	fileWriterMock := &writerTestdata.FileWriterMock{WriteErr: writeErr}
	cmdFactoryMock := testdata.CmdFactoryMock{GraphCmdName: "echo", ListCmdName: "echo"}
	job := NewJob("file", cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.ErrorIs(T, job.Error(), writeErr)
}

func TestRunCloseErr(T *testing.T) {
	closeErr := errors.New("close-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CloseErr: closeErr}
	cmdFactoryMock := testdata.CmdFactoryMock{GraphCmdName: "echo", ListCmdName: "echo"}
	job := NewJob("file", cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.ErrorIs(T, job.Error(), closeErr)
}

func TestRun(T *testing.T) {
	fileContents := []byte("MakeGraphCmd\n\nMakeListCmd\n")
	fileWriterMock := &writerTestdata.FileWriterMock{}
	cmdFactoryMock := testdata.CmdFactoryMock{GraphCmdName: "echo", ListCmdName: "echo"}
	job := NewJob("file", cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.NoError(T, job.Error())
	assert.Equal(T, fileContents, fileWriterMock.Contents)
}
