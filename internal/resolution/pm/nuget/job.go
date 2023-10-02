package nuget

import (
	"fmt"
	"os"

	"github.com/debricked/cli/internal/resolution/job"
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
		defer j.cleanupTempoCsproj()
		if err != nil {
			j.Errors().Critical(fmt.Errorf("%s\n%s", output, err))

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

func (j *Job) cleanupTempoCsproj() {
	// Cleanup of the temporary .csproj file (packages.config)
	tempFile := j.cmdFactory.GetTempoCsproj()
	if tempFile != "" {
		// remove the packages.config.csproj file
		err := osRemoveAll(tempFile)
		if err != nil {
			j.Errors().Critical(fmt.Errorf("failed to remove temporary .csproj file: %s", err))
		}
	}
}
