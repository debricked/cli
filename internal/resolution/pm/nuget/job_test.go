package nuget

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/nuget/testdata"
	"github.com/stretchr/testify/assert"
)

const (
	badName = "bad-name"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", false, CmdFactory{
		execPath: ExecPath{},
	})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestRunInstall(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", false, cmdFactoryMock)

	_, err := j.runInstallCmd()
	assert.NoError(t, err)

	assert.False(t, j.Errors().HasError())
}

func TestInstall(t *testing.T) {
	j := Job{install: true}
	assert.Equal(t, true, j.Install())

	j = Job{install: false}
	assert.Equal(t, false, j.Install())
}

func TestRunInstallCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeInstallErr = cmdErr
	j := NewJob("file", true, cmdFactoryMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), cmdErr)
}

func TestRunInstallCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.InstallCmdName = badName
	j := NewJob("file", true, cmdMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunInstallCmdOutputAndErr(t *testing.T) {
	// Setup mock to return a specific error and output
	expectedError := errors.New("expected error")
	expectedOutput := []byte("expected output")

	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.CmdOutput = expectedOutput
	cmdFactoryMock.CmdError = expectedError

	j := NewJob("file", true, cmdFactoryMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	// Check that the error recorded in j.Errors() matches the error created from the output
	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), errors.New(string(expectedOutput)))
}
