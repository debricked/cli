package pip

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/debricked/cli/pkg/resolution/pm/pip/testdata"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
	writerTestdata "github.com/debricked/cli/pkg/resolution/pm/writer/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	job := NewJob("file", false, CmdFactory{}, writer.FileWriter{})
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

// TODO add more tests a la maven / golang

func TestRun(t *testing.T) {
	// Load gt-data
	list, err := os.ReadFile("testdata/list.txt")
	assert.Nil(t, err)
	req, err := os.ReadFile("testdata/requirements.txt")
	assert.Nil(t, err)
	show, err := os.ReadFile("testdata/show.txt")
	assert.Nil(t, err)

	delimeter := "***"
	var fileContents []string
	fileContents = append(fileContents, string(req))
	fileContents = append(fileContents, delimeter)
	fileContents = append(fileContents, string(list))
	fileContents = append(fileContents, delimeter)
	fileContents = append(fileContents, string(show))
	res := []byte(strings.Join(fileContents, "\n"))

	fileWriterMock := &writerTestdata.FileWriterMock{}
	cmdFactoryMock := testdata.CmdFactoryMock{InstallCmdName: "echo", ListCmdName: "echo", CatCmdName: "echo", ShowCmdName: "echo"}
	job := NewJob("file", false, cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.NoError(t, job.Error())
	assert.Equal(t, res, fileWriterMock.Contents)
}
