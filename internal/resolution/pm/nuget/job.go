package nuget

import (
	"fmt"
	"os"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	nuget = "dotnet"
)

type Job struct {
	job.BaseJob
	install      bool
	nugetCommand string
	cmdFactory   ICmdFactory
}

func NewJob(
	file string,
	install bool,
	cmdFactory ICmdFactory,
) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(file),
		install:    install,
		cmdFactory: cmdFactory,
	}
}

func (j *Job) Install() bool {
	return j.install
}

func (j *Job) Run() {
	if j.install {

		j.SendStatus("installing dependencies")
		output, err := j.runInstallCmd()
		defer j.cleanupTempCsproj()
		if err != nil {
			formatted_error := fmt.Errorf("%s\n%s", output, err)
			j.Errors().Critical(util.NewPMJobError(formatted_error.Error()))

			return
		}
	}

}

var osRemoveAll = os.RemoveAll

func (j *Job) runInstallCmd() ([]byte, error) {

	j.nugetCommand = nuget
	installCmd, err := j.cmdFactory.MakeInstallCmd(j.nugetCommand, j.GetFile())
	if err != nil {
		return nil, err
	}

	installCmdOutput, err := installCmd.Output()
	if err != nil {
		return installCmdOutput, j.GetExitError(err)
	}

	return installCmdOutput, nil
}

func (j *Job) cleanupTempCsproj() {
	// Cleanup of the temporary .csproj file (packages.config)
	tempFile := j.cmdFactory.GetTempoCsproj()
	if tempFile != "" {
		// remove the packages.config.csproj file
		err := osRemoveAll(tempFile)
		formatted_error := fmt.Errorf("failed to remove temporary .csproj file: %s", err)
		if err != nil {
			j.Errors().Critical(util.NewPMJobError(formatted_error.Error()))
		}
	}
}
