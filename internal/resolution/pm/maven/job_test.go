package maven

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/maven/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", CmdFactory{})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestRunCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	j := NewJob("file", testdata.CmdFactoryMock{Err: cmdErr})

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), cmdErr)
}

func TestRunCmdOutputErr(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "bad-name"})

	go jobTestdata.WaitStatus(j)

	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunCmdOutputErrNoOutput(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "go", Arg: "bad-arg"})

	go jobTestdata.WaitStatus(j)

	j.Run()

	errs := j.Errors().GetAll()
	assert.Len(t, errs, 1)
	err := errs[0]

	// assert empty because, when Output is executed it will allocate memory for the byte slice to contain the standard output.
	// However since no bytes are sent to standard output err will be empty here.
	assert.Empty(t, err)
}

func TestRun(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "echo"})

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.False(t, j.Errors().HasError())
}
