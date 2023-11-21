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

	_, err := j.runInstallCmd()
	assert.NoError(t, err)

	assert.False(t, j.Errors().HasError())
}

func TestRunInstallPackagesConfig(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.GetTempoCsprojReturn = "tempo.csproj"
	j := NewJob("packages.config", false, cmdFactoryMock)

	_, err := j.runInstallCmd()
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

	go jobTestdata.WaitStatus(j)
	j.Run()
	errors := j.Errors().GetAll()
	assert.Equal(t, errors[0], util.NewPMJobError(cmdErrGt.Error()))

}

func TestInstall(t *testing.T) {
	j := Job{install: true}
	assert.Equal(t, true, j.Install())

	j = Job{install: false}
	assert.Equal(t, false, j.Install())
}

func TestRunInstallCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdErrGt := errors.New("\ncmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MakeInstallErr = cmdErr
	j := NewJob("file", true, cmdFactoryMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Equal(t, j.Errors().GetAll()[0], util.NewPMJobError(cmdErrGt.Error()))
}

func TestRunInstallCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.InstallCmdName = badName
	j := NewJob("file", true, cmdMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}
