package nuget

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/nuget/testdata"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/stretchr/testify/assert"
)

const (
	badName = "bad-name"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", false, &CmdFactory{
		execPath: ExecPath{},
	})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestRunInstall(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", false, cmdFactoryMock)

	_, _, err := j.runInstallCmd()
	assert.NoError(t, err)

	assert.False(t, j.Errors().HasError())
}

func TestRunInstallPackagesConfig(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.GetTempoCsprojReturn = "tempo.csproj"
	j := NewJob("packages.config", false, cmdFactoryMock)

	_, _, err := j.runInstallCmd()
	assert.NoError(t, err)

	assert.False(t, j.Errors().HasError())
}

func TestRunInstallPackagesConfigRemoveAllErr(t *testing.T) {

	oldOsRemoveAll := osRemoveAll
	cmdErr := errors.New("os-remove-all-error")
	cmdErrGt := errors.New("failed to remove temporary .csproj file: os-remove-all-error")
	osRemoveAll = func(path string) error {
		return cmdErr
	}

	defer func() {
		osRemoveAll = oldOsRemoveAll
	}()

	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.GetTempoCsprojReturn = "tempo.csproj"
	j := NewJob("packages.config", true, cmdFactoryMock)

	expectedError := util.NewPMJobError(cmdErrGt.Error())
	expectedError.SetStatus("cleanup")

	go jobTestdata.WaitStatus(j)
	j.Run()
	errors := j.Errors().GetAll()
	assert.Equal(t, errors[0], expectedError)

}

func TestInstall(t *testing.T) {
	j := Job{install: true}
	assert.Equal(t, true, j.Install())

	j = Job{install: false}
	assert.Equal(t, false, j.Install())
}

func TestRunInstallCmdErr(t *testing.T) {
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
			name:  "Invalid package version",
			error: "Unable to find a stable package 'PackageId' with version (>= 3.0.0)\n  - Found 10 version(s) in 'sourceA' [ Nearest version: '4.0.0-rc-2129' ]\n  - Found 9 version(s) in 'sourceB' [ Nearest version: '3.0.0-beta-00032' ]\n  - Found 0 version(s) in 'sourceC'\n  - Found 0 version(s) in 'sourceD'",
			doc:   "Failed to find package 'PackageId' with version (>= 3.0.0). Please check that package versions are correct in the manifest file.",
		},
		{
			name:  "Invalid package name",
			error: "Unable to find package 'PackageId'. No packages exist with this id in source(s): sourceA, sourceB, sourceC",
			doc:   "Failed to find package \"PackageId\" that satisfies the requirements. Please check that dependencies are correct in the manifest file. \nIf this is a private dependency, please make sure that the debricked CLI has access to install it or pre-install it before running the debricked CLI.",
		},
		{
			name:  "Unable to resolve package",
			error: "Unable to resolve 'Dependency (>= 1.0.0)' for 'TargetFramework'",
			doc:   "Couldn't resolve \"Dependency (>= 1.0.0)\". Please check if it exists and NuGet sources are configured properly.",
		},
		{
			name:  "No internet connection",
			error: "Unable to load the service index for source https://api.nuget.org/v3/index.json.\nAn error occurred while sending the request. \nThe remote name could not be resolved: 'api.nuget.org'",
			doc:   "Registry \"https://api.nuget.org/v3/index.json\" is not available at the moment. There might be a trouble with your network connection.",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cmdErr := errors.New(c.error)
			cmdFactoryMock := testdata.NewEchoCmdFactory()
			cmdFactoryMock.MakeInstallErr = cmdErr
			cmd, _ := cmdFactoryMock.MakeInstallCmd("echo", "package.json")

			expectedError := util.NewPMJobError("\n" + c.error)
			expectedError.SetDocumentation(c.doc)
			expectedError.SetStatus("installing dependencies")
			expectedError.SetCommand(cmd.String())

			j := NewJob("file", true, cmdFactoryMock)

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
	j := NewJob("file", true, cmdMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}
