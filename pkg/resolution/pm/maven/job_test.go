package maven

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/pkg/resolution/job/testdata"
	"github.com/debricked/cli/pkg/resolution/pm/maven/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", CmdFactory{})
	assert.Equal(t, "file", j.GetFile())
	assert.Nil(t, j.Error())
}

func TestRunCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	j := NewJob("file", testdata.CmdFactoryMock{Err: cmdErr})

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.ErrorIs(t, cmdErr, j.Error())
}

func TestRunCmdOutputErr(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "bad-name"})

	go jobTestdata.WaitStatus(j)

	j.Run()

	jobTestdata.AssertPathErr(t, j.Error())
}

func TestRun(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "echo"})

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.NoError(t, j.Error())
}
