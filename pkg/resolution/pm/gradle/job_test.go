package gradle

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/pkg/resolution/job/testdata"
	"github.com/debricked/cli/pkg/resolution/pm/gradle/testdata"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
	writerTestdata "github.com/debricked/cli/pkg/resolution/pm/writer/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", "nil", "nil", CmdFactory{}, writer.FileWriter{})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestRunCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	j := NewJob("file", "nil", "nil", testdata.CmdFactoryMock{Err: cmdErr}, writer.FileWriter{})

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), cmdErr)
}

func TestRunCmdOutputErr(t *testing.T) {
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: errors.New("create-error")}

	j := NewJob("file", "gradlew", "path", testdata.CmdFactoryMock{Name: "bad-name"}, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunCreateErr(t *testing.T) {
	createErr := errors.New("create-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}
	j := NewJob("file", "gradlew", "path", testdata.CmdFactoryMock{Name: "echo", Err: createErr}, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), createErr)
}

func TestRunWriteErr(t *testing.T) {
	writeErr := errors.New("write-error")
	fileWriterMock := &writerTestdata.FileWriterMock{WriteErr: writeErr}
	j := NewJob("file", "", "", testdata.CmdFactoryMock{Name: "echo", Err: writeErr}, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), writeErr)
}

func TestRunCloseErr(t *testing.T) {
	closeErr := errors.New("close-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CloseErr: closeErr}
	j := NewJob("file", "gradlew", "path", testdata.CmdFactoryMock{Name: "echo", Err: closeErr}, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), closeErr)
}

func TestRun(t *testing.T) {
	fileContents := []byte("MakeDependenciesCmd\n")
	fileWriterMock := &writerTestdata.FileWriterMock{Contents: fileContents}
	cmdFactoryMock := testdata.CmdFactoryMock{Name: "echo"}
	j := NewJob("file", "gradlew", "path", cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.False(t, j.Errors().HasError())
	assert.Equal(t, fileContents, fileWriterMock.Contents)
}
