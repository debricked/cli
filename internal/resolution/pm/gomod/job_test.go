package gomod

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/gomod/testdata"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/debricked/cli/internal/resolution/pm/writer"
	writerTestdata "github.com/debricked/cli/internal/resolution/pm/writer/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", CmdFactory{}, writer.FileWriter{})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestRunGraphCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeGraphCmdErr = cmdErr
	j := NewJob("file", cmdFactoryMock, nil)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Contains(t, j.Errors().GetCriticalErrors(), util.NewPMJobError(cmdErr.Error()))
}

func TestRunCmdOutputErr(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.GraphCmdName = "bad-name"
	j := NewJob("file", cmdFactoryMock, nil)

	go jobTestdata.WaitStatus(j)

	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunListCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeListCmdErr = cmdErr
	j := NewJob("file", cmdFactoryMock, nil)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), util.NewPMJobError(cmdErr.Error()))
}

func TestRunListCmdOutputErr(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.ListCmdName = "bad-name"
	j := NewJob("file", cmdFactoryMock, nil)

	go jobTestdata.WaitStatus(j)

	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunCreateErr(t *testing.T) {
	createErr := errors.New("create-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), util.NewPMJobError(createErr.Error()))
}

func TestRunWriteErr(t *testing.T) {
	writeErr := errors.New("write-error")
	fileWriterMock := &writerTestdata.FileWriterMock{WriteErr: writeErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), util.NewPMJobError(writeErr.Error()))
}

func TestRunCloseErr(t *testing.T) {
	closeErr := errors.New("close-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CloseErr: closeErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), util.NewPMJobError(closeErr.Error()))
}

func TestRun(t *testing.T) {
	fileContents := []byte("MakeGraphCmd\n\nMakeListCmd\n")
	fileWriterMock := &writerTestdata.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Empty(t, j.Errors().GetAll())
	assert.Equal(t, fileContents, fileWriterMock.Contents)
}
