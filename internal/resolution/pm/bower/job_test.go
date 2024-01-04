package bower

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/bower/testdata"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/debricked/cli/internal/resolution/pm/writer"
	writerTestdata "github.com/debricked/cli/internal/resolution/pm/writer/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", CmdFactory{}, writer.FileWriter{})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestRunGraphCmdErr(t *testing.T) {
	cases := []struct {
		name  string
		error string
		doc   string
	}{
		{
			name:  "General error",
			error: "cmd-error",
			doc:   util.UnknownError,
		},
		{
			name:  "Bower not found",
			error: "        |exec: \"bower\": executable file not found in $PATH",
			doc:   "Bower wasn't found. Please check if it is installed and accessible by the CLI.",
		},
		{
			name:  "Invalid package version",
			error: "bower get-size#11.0.*       not-cached https://github.com/desandro/get-size.git#11.0.*\nbower get-size#11.0.*          resolve https://github.com/desandro/get-size.git#11.0.*\nbower get-size#11.0.*     ENORESTARGET No tag found that was able to satisfy 11.0.*\n\nAdditional error details:\nAvailable versions in https://github.com/desandro/get-size.git: 3.0.0, 2.0.3, 2.0.2",
			doc:   "Failed to find package \"get-size#11.0.*\" that satisfies the requirement from bower.json file. In most cases you or one of your dependencies are requesting a package version that doesn't exist. Please check that package versions are correct.",
		},
		{
			name:  "Invalid package name",
			error: "bower get-sizee#3.0.*        ENOTFOUND Package get-sizee not foundbower get-sizee#3.0.*        ENOTFOUND Package get-sizee not found",
			doc:   "Failed to find package \"get-sizee\" that satisfies the requirements. Please check that dependencies are correct in your bower.json file. \nIf this is a private dependency, please make sure that the debricked CLI has access to install it or pre-install it before running the debricked CLI.",
		},
		{
			name:  "Permission denied",
			error: "bower ev-emitter#^2.0.0     not-cached https://github.com/metafizzy/ev-emitter.git#^2.0.0\nbower ev-emitter#^2.0.0        resolve https://github.com/metafizzy/ev-emitter.git#^2.0.0\nbower ev-emitter#^2.0.0       download https://github.com/metafizzy/ev-emitter/archive/v2.1.2.tar.gz\nbower ev-emitter#^2.0.0        extract archive.tar.gz\nbower ev-emitter#^2.0.0       resolved https://github.com/metafizzy/ev-emitter.git#2.1.2\nbower ev-emitter#^2.0.0        install ev-emitter#2.1.2\nbower                           EACCES EACCES: permission denied, rmdir '/home/asus/Projects/playground/bowerr/bower_components/ev-emitter'\n\nStack trace:\nError: EACCES: permission denied, rmdir '/home/asus/Projects/playground/bowerr/bower_components/ev-emitter'",
			doc:   "Couldn't get access to \"/home/asus/Projects/playground/bowerr/bower_components/ev-emitter\". Please check permissions or try running this command again as root/Administrator.",
		},
		{
			name:  "No internet connection",
			error: "bower ev-emitterr#^1.0.0     EAI_AGAIN Request to https://registry.bower.io/packages/ev-emitterr failed: getaddrinfo EAI_AGAIN registry.bower.io",
			doc:   "Package registry \"registry.bower.io\" is not available at the moment. There might be a trouble with your network connection.",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cmdErr := errors.New(c.error)
			cmdFactoryMock := testdata.NewEchoCmdFactory()
			cmdFactoryMock.MakeInstallCmdErr = cmdErr
			cmd, _ := cmdFactoryMock.MakeInstallCmd("echo", "")

			expectedError := util.NewPMJobError(c.error)
			expectedError.SetDocumentation(c.doc)
			expectedError.SetStatus("installing dependencies")
			expectedError.SetCommand(cmd.String())

			j := NewJob("file", cmdFactoryMock, nil)

			go jobTestdata.WaitStatus(j)
			j.Run()

			allErrors := j.Errors().GetAll()

			assert.Len(t, j.Errors().GetAll(), 1)
			assert.Contains(t, allErrors, expectedError)
		})
	}
}

func TestRunCmdOutputErr(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.InstallCmdName = "bad-name"
	j := NewJob("file", cmdFactoryMock, nil)

	go jobTestdata.WaitStatus(j)

	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunListCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeListCmdErr = cmdErr
	j := NewJob("file", cmdFactoryMock, nil)

	cmd, _ := cmdFactoryMock.MakeListCmd("echo", "")
	expectedError := util.NewPMJobError(cmdErr.Error())
	expectedError.SetStatus("creating dependency tree")
	expectedError.SetCommand(cmd.String())

	go jobTestdata.WaitStatus(j)

	j.Run()

	allErrors := j.Errors().GetAll()

	assert.Len(t, allErrors, 1)
	assert.Contains(t, allErrors, expectedError)
}

func TestRunListCmdOutputErr(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.ListCmdName = "bad-name"
	j := NewJob("file", cmdFactoryMock, nil)

	go jobTestdata.WaitStatus(j)

	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunCreateErr(t *testing.T) {
	createErr := errors.New("create-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	cmd, _ := cmdFactoryMock.MakeListCmd("echo", "")
	expectedError := util.NewPMJobError(createErr.Error())
	expectedError.SetStatus("creating lock file")
	expectedError.SetCommand(cmd.String())

	go jobTestdata.WaitStatus(j)

	j.Run()

	allErrors := j.Errors().GetAll()

	assert.Len(t, allErrors, 1)
	assert.Contains(t, allErrors, expectedError)
}

func TestRunWriteErr(t *testing.T) {
	writeErr := errors.New("write-error")
	fileWriterMock := &writerTestdata.FileWriterMock{WriteErr: writeErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	cmd, _ := cmdFactoryMock.MakeListCmd("echo", "")
	expectedError := util.NewPMJobError(writeErr.Error())
	expectedError.SetStatus("creating lock file")
	expectedError.SetCommand(cmd.String())

	go jobTestdata.WaitStatus(j)

	j.Run()

	allErrors := j.Errors().GetAll()

	assert.Len(t, allErrors, 1)
	assert.Contains(t, allErrors, expectedError)
}

func TestRunCloseErr(t *testing.T) {
	closeErr := errors.New("close-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CloseErr: closeErr}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	allErrors := j.Errors().GetAll()

	assert.Len(t, allErrors, 1)
	assert.Contains(t, allErrors, util.NewPMJobError(closeErr.Error()))
}

func TestRun(t *testing.T) {
	fileWriterMock := &writerTestdata.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Empty(t, j.Errors().GetAll())
	assert.Equal(t, "MakeListCmd\n", string(fileWriterMock.Contents))
}
