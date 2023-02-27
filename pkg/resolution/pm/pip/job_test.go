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

const (
	badName = "bad-name"
)

func TestNewJob(t *testing.T) {
	job := NewJob("file", false, CmdFactory{
		execPath: ExecPath{},
	}, writer.FileWriter{})
	assert.Equal(t, "file", job.file)
	assert.Nil(t, job.err)
}

func TestFile(t *testing.T) {
	job := Job{file: "file"}
	assert.Equal(t, "file", job.File())
}

func TestInstall(t *testing.T) {
	job := Job{install: true}
	assert.Equal(t, true, job.Install())

	job = Job{install: false}
	assert.Equal(t, false, job.Install())
}

func TestError(t *testing.T) {
	jobErr := errors.New("error")
	job := Job{file: "file", err: jobErr}
	assert.Equal(t, jobErr, job.Error())
}

func TestRunCreateVenvCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeCreateVenvErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	job := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.ErrorIs(t, cmdErr, job.Error())
}

func TestRunCreateVenvCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.CreateVenvCmdName = badName
	job := NewJob("file", true, cmdMock, nil)
	job.Run()
	assert.ErrorContains(t, job.err, "executable file not found in")
	assert.ErrorContains(t, job.err, "PATH")
}

func TestRunInstallCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeInstallErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	job := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.ErrorIs(t, cmdErr, job.Error())
}

func TestRunInstallCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.InstallCmdName = badName
	job := NewJob("file", true, cmdMock, nil)
	job.Run()
	assert.ErrorContains(t, job.err, "executable file not found in")
	assert.ErrorContains(t, job.err, "PATH")
}

func TestRunCatCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeCatErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	job := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.ErrorIs(t, cmdErr, job.Error())
}

func TestRunCatCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.CatCmdName = badName
	job := NewJob("file", false, cmdMock, nil)
	job.Run()
	assert.ErrorContains(t, job.err, "executable file not found in")
	assert.ErrorContains(t, job.err, "PATH")
}

func TestRunListCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeListErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	job := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.ErrorIs(t, cmdErr, job.Error())
}

func TestRunListCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.ListCmdName = badName
	job := NewJob("file", false, cmdMock, nil)
	job.Run()
	assert.ErrorContains(t, job.err, "executable file not found in")
	assert.ErrorContains(t, job.err, "PATH")
}

func TestRunShowCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeShowErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	job := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	job.Run()

	assert.ErrorIs(t, cmdErr, job.Error())
}

func TestRunShowCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.ShowCmdName = badName
	job := NewJob("file", false, cmdMock, nil)
	job.Run()
	assert.ErrorContains(t, job.err, "executable file not found in")
	assert.ErrorContains(t, job.err, "PATH")
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

func TestRunInstall(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	fileWriterMock := &writerTestdata.FileWriterMock{}
	job := NewJob("file", false, cmdFactoryMock, fileWriterMock)

	_, err := job.runInstallCmd()
	assert.NoError(t, err)

	assert.NoError(t, job.Error())
}

func TestParsePipList(t *testing.T) {
	job := NewJob("file", false, CmdFactory{
		execPath: ExecPath{},
	}, writer.FileWriter{})
	file, err := os.ReadFile("testdata/list.txt")
	assert.Nil(t, err)
	pipData := string(file)
	packages := job.parsePipList(pipData)
	gt := []string{"aiohttp", "cryptography", "numpy", "Flask", "open-source-health", "pandas", "tqdm"}
	assert.Equal(t, gt, packages)
	assert.Nil(t, job.err)
}

func TestRunCreateErr(t *testing.T) {
	createErr := errors.New("create-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}
	cmdMock := testdata.NewEchoCmdFactory()
	job := NewJob("file", true, cmdMock, fileWriterMock)
	job.Run()

	assert.ErrorIs(t, job.Error(), createErr)
}

func TestRunWriteErr(t *testing.T) {
	writeErr := errors.New("write-error")
	fileWriterMock := &writerTestdata.FileWriterMock{WriteErr: writeErr}
	cmdMock := testdata.NewEchoCmdFactory()
	job := NewJob("file", true, cmdMock, fileWriterMock)
	job.Run()

	assert.ErrorIs(t, job.Error(), writeErr)
}

func TestRunCloseErr(t *testing.T) {
	closeErr := errors.New("close-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CloseErr: closeErr}
	cmdMock := testdata.NewEchoCmdFactory()
	job := NewJob("file", true, cmdMock, fileWriterMock)
	job.Run()

	assert.ErrorIs(t, job.Error(), closeErr)
}
