package poetry

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/poetry/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestRunCmdErrExecutableNotFound(t *testing.T) {
	execErr := errors.New("exec: \"poetry\": executable file not found in $PATH")
	j := NewJob("file", testdata.CmdFactoryMock{Err: execErr})

	go jobTestdata.WaitStatus(j)

	j.Run()

	errs := j.Errors().GetAll()
	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "executable file not found")
	assert.Contains(t, errs[0].Documentation(), "Poetry wasn't found")
}

func TestRunSuccess(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "echo", Arg: "ok"})

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.False(t, j.Errors().HasError())
	assert.Len(t, j.Errors().GetAll(), 0)
}
