package maven

import (
	"errors"
	"fmt"
	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/pm/maven/testdata"
	"runtime"

	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", CmdFactory{})
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

func TestRunCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	j := NewJob("file", testdata.CmdFactoryMock{Err: cmdErr})

	go waitStatus(j)

	j.Run()

	assert.ErrorIs(t, cmdErr, j.Error())
}

func TestRunCmdOutputErr(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "bad-name"})

	go waitStatus(j)

	j.Run()

	assertPathErr(t, j.Error())
}

func TestRun(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "echo"})

	go waitStatus(j)

	j.Run()

	assert.NoError(t, j.err)
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
