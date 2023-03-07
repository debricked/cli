package gomod

import (
	"errors"
	"fmt"
	"runtime"
	"testing"

	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/pm/gomod/testdata"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
	writerTestdata "github.com/debricked/cli/pkg/resolution/pm/writer/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", CmdFactory{}, writer.FileWriter{})
	assert.Equal(t, "file", j.file)
	assert.Nil(t, j.err)
}

func TestFile(t *testing.T) {
	j := Job{file: "file"}
	assert.Equal(t, "file", j.File())
}

func TestError(t *testing.T) {
	jobErr := errors.New("error")
	j := Job{file: "file", err: jobErr}
	assert.Equal(t, jobErr, j.Error())
}

func TestRunGraphCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeGraphCmdErr = cmdErr
	j := NewJob("file", cmdFactoryMock, nil)

	go waitStatus(j)

	j.Run()

	assert.ErrorIs(t, cmdErr, j.Error())
}

func TestRunCmdOutputErr(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.GraphCmdName = "bad-name"
	j := NewJob("file", cmdFactoryMock, nil)

	go waitStatus(j)

	j.Run()

	assertPathErr(t, j.Error())
}

func TestRunListCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeListCmdErr = cmdErr
	j := NewJob("file", cmdFactoryMock, nil)

	go waitStatus(j)

	j.Run()

	assert.ErrorIs(t, cmdErr, j.Error())
}

func TestRunListCmdOutputErr(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.ListCmdName = "bad-name"
	j := NewJob("file", cmdFactoryMock, nil)

	go waitStatus(j)

	j.Run()

	assertPathErr(t, j.Error())
}

func TestRunCreateErr(t *testing.T) {
	createErr := errors.New("create-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	go waitStatus(j)

	j.Run()

	assert.ErrorIs(t, j.Error(), createErr)
}

func TestRunWriteErr(t *testing.T) {
	writeErr := errors.New("write-error")
	fileWriterMock := &writerTestdata.FileWriterMock{WriteErr: writeErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	go waitStatus(j)

	j.Run()

	assert.ErrorIs(t, j.Error(), writeErr)
}

func TestRunCloseErr(t *testing.T) {
	closeErr := errors.New("close-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CloseErr: closeErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	go waitStatus(j)

	j.Run()

	assert.ErrorIs(t, j.Error(), closeErr)
}

func TestRun(t *testing.T) {
	fileContents := []byte("MakeGraphCmd\n\nMakeListCmd\n")
	fileWriterMock := &writerTestdata.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	go waitStatus(j)

	j.Run()

	assert.NoError(t, j.Error())
	assert.Equal(t, fileContents, fileWriterMock.Contents)
}

func waitStatus(j job.IJob) {
	for {
		<-j.Status()
	}
}

func assertPathErr(t *testing.T, err error) {
	var path string
	if runtime.GOOS == "windows" {
		path = "%PATH%"
	} else {
		path = "$PATH"
	}
	errMsg := fmt.Sprintf("executable file not found in %s", path)
	assert.ErrorContains(t, err, errMsg)
}
