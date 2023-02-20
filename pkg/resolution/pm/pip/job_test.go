package pip

import (
	"errors"
	"fmt"
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
func TestRunCreateVenvCmdErr(T *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeCreateVenvErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	job := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.ErrorIs(T, cmdErr, job.Error())
}

func TestRunActivateVenvCmdErr(T *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeActivateVenvErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	job := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.ErrorIs(T, cmdErr, job.Error())
}

func TestRunInstallCmdErr(T *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeInstallErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	job := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.ErrorIs(T, cmdErr, job.Error())
}

func TestRunCatCmdErr(T *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeCatErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	job := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.ErrorIs(T, cmdErr, job.Error())
}

func TestRunListCmdErr(T *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeListErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	job := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.ErrorIs(T, cmdErr, job.Error())
}

func TestRunShowCmdErr(T *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeShowErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	job := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.ErrorIs(T, cmdErr, job.Error())
}

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
	fileContents = append(fileContents, string(req)+"\n")
	fileContents = append(fileContents, delimeter)
	fileContents = append(fileContents, string(list)+"\n")
	fileContents = append(fileContents, delimeter)
	fileContents = append(fileContents, string(show)+"\n")
	res := []byte(strings.Join(fileContents, "\n"))

	fileWriterMock := &writerTestdata.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	job := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.NoError(t, job.Error())
	fmt.Println(string(fileWriterMock.Contents))
	assert.Equal(t, string(res), string(fileWriterMock.Contents))
}

func TestParsePipList(t *testing.T) {
	job := NewJob("file", false, CmdFactory{}, writer.FileWriter{})
	file, err := os.ReadFile("testdata/list.txt")
	assert.Nil(t, err)
	pipData := string(file)
	packages, err := job.parsePipList(pipData)
	assert.Nil(t, err)
	gt := []string{"aiohttp", "cryptography", "numpy", "Flask", "open-source-health", "pandas", "tqdm"}
	assert.Equal(t, gt, packages)
	assert.Nil(t, job.err)
}
