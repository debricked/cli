package java

import (
	"errors"
	"fmt"
	"testing"

	conf "github.com/debricked/cli/pkg/callgraph/config"
	jobTestdata "github.com/debricked/cli/pkg/callgraph/job/testdata"
	"github.com/debricked/cli/pkg/callgraph/language/java11/testdata"
	ioWriter "github.com/debricked/cli/pkg/io/writer"
	writerTestdata "github.com/debricked/cli/pkg/io/writer/testdata"
	"github.com/stretchr/testify/assert"
)

const (
	badName = "bad-name"
	dir     = "dir"
)

var files = []string{"file"}

func TestNewJob(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	writer := ioWriter.FileWriter{}
	config := conf.Config{}
	j := NewJob(dir, files, cmdFactoryMock, writer, config)
	assert.Equal(t, "file", j.GetFiles())
	assert.Equal(t, "dir", j.GetDir())
	assert.False(t, j.Errors().HasError())
}

// func TestRunMakeGradleCopyDependenciesCmdErr(t *testing.T) {
// 	cmdErr := errors.New("cmd-error")
// 	cmdFactoryMock := testdata.NewEchoCmdFactory()
// 	cmdFactoryMock.GradleCopyDepErr = cmdErr
// 	fileWriterMock := &writerTestdata.FileWriterMock{}
// 	config := conf.NewConfig("java", nil, map[string]string{"pm": gradle})
// 	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, config)

// 	go jobTestdata.WaitStatus(j)
// 	j.Run()

// 	assert.Len(t, j.Errors().GetAll(), 1)
// 	assert.Contains(t, j.Errors().GetAll(), cmdErr)
// }

// func TestRunMakeGradleCopyDependenciesOutputErr(t *testing.T) {
// 	cmdMock := testdata.NewEchoCmdFactory()
// 	cmdMock.GradleCopyDepName = badName
// 	cmdFactoryMock := testdata.NewEchoCmdFactory()
// 	fileWriterMock := &writerTestdata.FileWriterMock{}
// 	config := conf.NewConfig("java", nil, map[string]string{"pm": gradle})
// 	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, config)

// 	go jobTestdata.WaitStatus(j)
// 	j.Run()

// 	jobTestdata.AssertPathErr(t, j.Errors())
// }

func TestRunMakeMavenCopyDependenciesCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MvnCopyDepErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven})
	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, config)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), cmdErr)
}

func TestRunMakeMavenCopyDependenciesOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.MvnCopyDepName = badName
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	fileWriterMock := &writerTestdata.FileWriterMock{}
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven})
	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, config)

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRun(t *testing.T) {

	fileWriterMock := &writerTestdata.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven})
	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, config)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.False(t, j.Errors().HasError())
	fmt.Println(string(fileWriterMock.Contents))
	assert.False(t, false)
}

func TestRunCreateErr(t *testing.T) {
	createErr := errors.New("create-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven})
	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, config)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), createErr)
}

func TestRunWriteErr(t *testing.T) {
	writeErr := errors.New("write-error")
	fileWriterMock := &writerTestdata.FileWriterMock{WriteErr: writeErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven})
	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, config)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), writeErr)
}

func TestRunCloseErr(t *testing.T) {
	closeErr := errors.New("close-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CloseErr: closeErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven})
	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, config)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), closeErr)
}
