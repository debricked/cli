package pip

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	jobTestdata "github.com/debricked/cli/pkg/resolution/job/testdata"
	"github.com/debricked/cli/pkg/resolution/pm/pip/testdata"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
	writerTestdata "github.com/debricked/cli/pkg/resolution/pm/writer/testdata"
	"github.com/stretchr/testify/assert"
)

const (
	badName = "bad-name"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", false, CmdFactory{
		execPath: ExecPath{},
	}, writer.FileWriter{})
	assert.Equal(t, "file", j.File)
	assert.Nil(t, j.Err)
}

func TestInstall(t *testing.T) {
	j := Job{install: true}
	assert.Equal(t, true, j.Install())

	j = Job{install: false}
	assert.Equal(t, false, j.Install())
}

func TestRunCreateVenvCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeCreateVenvErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	j := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.ErrorIs(t, cmdErr, j.Error())
}

func TestRunCreateVenvCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.CreateVenvCmdName = badName
	j := NewJob("file", true, cmdMock, nil)

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Error())
}

func TestRunInstallCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeInstallErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	j := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.ErrorIs(t, cmdErr, j.Error())
}

func TestRunInstallCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.InstallCmdName = badName
	j := NewJob("file", true, cmdMock, nil)

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Error())
}

func TestRunCatCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeCatErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	j := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.ErrorIs(t, cmdErr, j.Error())
}

func TestRunCatCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.CatCmdName = badName
	j := NewJob("file", false, cmdMock, nil)

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Error())
}

func TestRunListCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeListErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	j := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.ErrorIs(t, cmdErr, j.Error())
}

func TestRunListCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.ListCmdName = badName
	j := NewJob("file", false, cmdMock, nil)

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Error())
}

func TestRunShowCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeShowErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	j := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.ErrorIs(t, cmdErr, j.Error())
}

func TestRunShowCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.ShowCmdName = badName
	j := NewJob("file", false, cmdMock, nil)

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Error())
}

func TestRun(t *testing.T) {
	// Load gt-data
	list, err := os.ReadFile("testdata/list.txt")
	assert.Nil(t, err)
	req, err := os.ReadFile("testdata/requirements.txt")
	assert.Nil(t, err)
	show, err := os.ReadFile("testdata/show.txt")
	assert.Nil(t, err)

	var fileContents []string
	fileContents = append(fileContents, string(req)+"\n")
	fileContents = append(fileContents, lockFileDelimiter)
	fileContents = append(fileContents, string(list)+"\n")
	fileContents = append(fileContents, lockFileDelimiter)
	fileContents = append(fileContents, string(show)+"\n")
	res := []byte(strings.Join(fileContents, "\n"))

	fileWriterMock := &writerTestdata.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", true, cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.NoError(t, j.Error())
	fmt.Println(string(fileWriterMock.Contents))
	assert.Equal(t, string(res), string(fileWriterMock.Contents))
}

func TestRunInstall(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	fileWriterMock := &writerTestdata.FileWriterMock{}
	j := NewJob("file", false, cmdFactoryMock, fileWriterMock)

	_, err := j.runInstallCmd()
	assert.NoError(t, err)

	assert.NoError(t, j.Error())
}

func TestParsePipList(t *testing.T) {
	j := NewJob("file", false, CmdFactory{
		execPath: ExecPath{},
	}, writer.FileWriter{})
	file, err := os.ReadFile("testdata/list.txt")
	assert.Nil(t, err)
	pipData := string(file)
	packages := j.parsePipList(pipData)
	gt := []string{"aiohttp", "cryptography", "numpy", "Flask", "open-source-health", "pandas", "tqdm"}
	assert.Equal(t, gt, packages)
	assert.Nil(t, j.Err)
}

func TestRunCreateErr(t *testing.T) {
	createErr := errors.New("create-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}
	cmdMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", true, cmdMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.ErrorIs(t, j.Error(), createErr)
}

func TestRunWriteErr(t *testing.T) {
	writeErr := errors.New("write-error")
	fileWriterMock := &writerTestdata.FileWriterMock{WriteErr: writeErr}
	cmdMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", true, cmdMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.ErrorIs(t, j.Error(), writeErr)
}

func TestRunCloseErr(t *testing.T) {
	closeErr := errors.New("close-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CloseErr: closeErr}
	cmdMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", true, cmdMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.ErrorIs(t, j.Error(), closeErr)
}
