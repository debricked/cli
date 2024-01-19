package pip

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/pip/testdata"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/debricked/cli/internal/resolution/pm/writer"
	writerTestdata "github.com/debricked/cli/internal/resolution/pm/writer/testdata"
	"github.com/stretchr/testify/assert"
)

const (
	badName = "bad-name"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", false, CmdFactory{
		execPath: ExecPath{},
	}, writer.FileWriter{}, pipCleaner{})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestInstall(t *testing.T) {
	j := Job{install: true}
	assert.Equal(t, true, j.Install())

	j = Job{install: false}
	assert.Equal(t, false, j.Install())
}

func TestRunCreateVenvCmdErr(t *testing.T) {
	cases := []struct {
		name  string
		error string
		doc   string
	}{
		{
			name:  "General error",
			error: "some random error",
			doc:   "Error when trying to create python virtual environment",
		},
		{
			name:  "Python not found",
			error: "        |exec: \"python\": executable file not found in $PATH",
			doc:   "Python wasn't found. Please check if it is installed and accessible by the CLI.",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cmdErr := errors.New(c.error)
			cmdFactoryMock := testdata.NewEchoCmdFactory()
			cmdFactoryMock.MakeCreateVenvErr = cmdErr
			cmd, _ := cmdFactoryMock.MakeCreateVenvCmd("file.venv")
			fileWriterMock := &writerTestdata.FileWriterMock{}
			j := NewJob("file", true, cmdFactoryMock, fileWriterMock, pipCleaner{})

			expectedError := util.NewPMJobError(c.error)
			expectedError.SetDocumentation(c.doc)
			expectedError.SetStatus("creating venv")
			expectedError.SetCommand(cmd.String())

			go jobTestdata.WaitStatus(j)
			j.Run()

			allErrors := j.Errors().GetAll()

			assert.Len(t, allErrors, 1)
			assert.Contains(t, allErrors, expectedError)
		})
	}
}

func TestRunCreateVenvCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.CreateVenvCmdName = badName
	j := NewJob("file", true, cmdMock, nil, pipCleaner{})

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunInstallCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeInstallErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	j := NewJob("file", true, cmdFactoryMock, fileWriterMock, pipCleaner{})

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
}

func TestRunInstallCmdErrors(t *testing.T) {
	cases := []struct {
		name  string
		error string
		doc   string
	}{
		{
			name:  "Pip not found",
			error: "        |exec: \"pip\": executable file not found in $PATH",
			doc:   "Pip wasn't found. Please check if it is installed and accessible by the CLI.",
		},
		{
			name:  "Python not found",
			error: "        |exec: \"python\": executable file not found in $PATH",
			doc:   "Python wasn't found. Please check if it is installed and accessible by the CLI.",
		},
		{
			name:  "Python3 not found",
			error: "        |exec: \"python3\": executable file not found in $PATH",
			doc:   "Python3 wasn't found. Please check if it is installed and accessible by the CLI.",
		},
		{
			name:  "Build Error",
			error: " python setup.py bdist_wheel did not run successfully. ",
			doc:   "Failed to build python dependency \"bdist_wheel\" with setup.py. This probably means the project was not set up correctly and could mean that an OS package is missing.",
		},
		{
			name:  "Auth Error",
			error: "WARNING: 401 Error, Credentials not correct for <some-pip-registry>\nNo matching distribution found for some-dependency>=0.1.3\n",
			doc:   "Failed to install python dependency \"some-dependency>=0.1.3\" due to authorization.\nIf this is a private dependency, please make sure that the debricked CLI has access to install it or pre-install it before running the debricked CLI.",
		},
		{
			name:  "Version Error",
			error: "Could not find a version that satisfies the requirement test==123",
			doc:   "Failed to find a version that satisfies the requirement for python dependency \"test\". This could mean that the package or version does not exist.\nIf this is a private dependency, please make sure that the debricked CLI has access to install it or pre-install it before running the debricked CLI.",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cmdErr := errors.New(c.error)
			cmdFactoryMock := testdata.NewEchoCmdFactory()
			cmdFactoryMock.MakeInstallErr = cmdErr
			cmd, _ := cmdFactoryMock.MakeInstallCmd("echo", "file")
			fileWriterMock := &writerTestdata.FileWriterMock{}

			expectedError := util.NewPMJobError(c.error)
			expectedError.SetDocumentation(c.doc)
			expectedError.SetStatus("installing dependencies")
			expectedError.SetCommand(cmd.String())

			j := NewJob("file", true, cmdFactoryMock, fileWriterMock, pipCleaner{})

			go jobTestdata.WaitStatus(j)
			j.Run()

			allErrors := j.Errors().GetAll()

			assert.Len(t, j.Errors().GetAll(), 1)
			assert.Contains(t, allErrors, expectedError)
		})
	}
}

func TestRunInstallCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.InstallCmdName = badName
	j := NewJob("file", true, cmdMock, nil, pipCleaner{})

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunCatCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeCatErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	j := NewJob("file", true, cmdFactoryMock, fileWriterMock, pipCleaner{})

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
}

func TestRunCatCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.CatCmdName = badName
	j := NewJob("file", false, cmdMock, nil, nil)

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunListCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeListErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	j := NewJob("file", true, cmdFactoryMock, fileWriterMock, pipCleaner{})

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
}

func TestRunListCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.ListCmdName = badName
	j := NewJob("file", false, cmdMock, nil, nil)

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunShowCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeShowErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}
	j := NewJob("file", true, cmdFactoryMock, fileWriterMock, pipCleaner{})

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
}

func TestRunShowCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.ShowCmdName = badName
	j := NewJob("file", false, cmdMock, nil, nil)

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
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
	j := NewJob("file", true, cmdFactoryMock, fileWriterMock, pipCleaner{})

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.False(t, j.Errors().HasError())
	fmt.Println(string(fileWriterMock.Contents))
	assert.Equal(t, string(res), string(fileWriterMock.Contents))
}

func TestRunInstall(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	fileWriterMock := &writerTestdata.FileWriterMock{}
	j := NewJob("file", false, cmdFactoryMock, fileWriterMock, nil)

	_, err := j.runInstallCmd()
	assert.NoError(t, err)

	assert.False(t, j.Errors().HasError())
}

func TestRunInstallWVenvPath(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	fileWriterMock := &writerTestdata.FileWriterMock{}
	j := NewJob("file", false, cmdFactoryMock, fileWriterMock, nil)
	j.venvPath = "test-path"

	_, err := j.runInstallCmd()

	var expectedPath string
	if runtime.GOOS == "windows" {
		expectedPath = "test-path/Scripts/pip"
	} else {
		expectedPath = "test-path/bin/pip"
	}

	expectedPath = filepath.FromSlash(expectedPath)
	assert.Equal(t, expectedPath, j.pipCommand)
	assert.NoError(t, err)

	assert.False(t, j.Errors().HasError())
}

func TestParsePipList(t *testing.T) {
	j := NewJob("file", false, CmdFactory{
		execPath: ExecPath{},
	}, writer.FileWriter{}, pipCleaner{})
	file, err := os.ReadFile("testdata/list.txt")
	assert.Nil(t, err)
	pipData := string(file)
	packages := j.parsePipList(pipData)
	gt := []string{"aiohttp", "cryptography", "numpy", "Flask", "open-source-health", "pandas", "tqdm"}
	assert.Equal(t, gt, packages)
	assert.False(t, j.Errors().HasError())
}

func TestRunCreateErr(t *testing.T) {
	createErr := errors.New("create-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}
	cmdMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", true, cmdMock, fileWriterMock, pipCleaner{})

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
}

func TestRunWriteErr(t *testing.T) {
	writeErr := errors.New("write-error")
	fileWriterMock := &writerTestdata.FileWriterMock{WriteErr: writeErr}
	cmdMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", true, cmdMock, fileWriterMock, pipCleaner{})

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
}

func TestRunCloseErr(t *testing.T) {
	closeErr := errors.New("close-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CloseErr: closeErr}
	cmdMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", true, cmdMock, fileWriterMock, pipCleaner{})

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	correctErr := util.NewPMJobError(closeErr.Error())
	assert.Contains(t, j.Errors().GetAll(), correctErr)
}

type pipCleanerMock struct {
	CleanErr error
}

func (p *pipCleanerMock) RemoveAll(_ string) error {
	return p.CleanErr
}

func TestRunCleanErr(t *testing.T) {
	CleanErr := errors.New("clean-error")
	fileWriterMock := &writerTestdata.FileWriterMock{}
	cmdMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", true, cmdMock, fileWriterMock, nil)
	j.pipCleaner = &pipCleanerMock{CleanErr: CleanErr}

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
}

var wasCalled bool

type pipCleanerMockCalled struct {
	WasCalled bool
}

func (p pipCleanerMockCalled) RemoveAll(_ string) error {
	wasCalled = true

	return nil
}

func TestErrorStillClean(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeInstallErr = cmdErr
	fileWriterMock := &writerTestdata.FileWriterMock{}

	wasCalled = false
	cleaner := pipCleanerMockCalled{}
	j := NewJob("file", true, cmdFactoryMock, fileWriterMock, cleaner)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.True(t, wasCalled)
}
