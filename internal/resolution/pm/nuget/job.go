package nuget

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	nuget                      = "dotnet"
	executableNotFoundErrRegex = `executable file not found`
	versionNotFoundErrRegex    = `Unable to find [\w\s]*package ('[^"'\s:]+' with version [^"'\n:]+)`
	dependencyNotFoundErrRegex = `'([^"'\s:]+)'. No packages exist`
	unableToResolveErrRegex    = `Unable to resolve '([^"'\n:]+)'`
	noInternetErrRegex         = `Unable to load the service index for source ([^"'\s]+).`
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
		status := "installing dependencies"
		j.SendStatus(status)
		output, cmd, err := j.runInstallCmd()
		defer j.cleanupTempCsproj()
		if err != nil {
			formatted_error := fmt.Errorf("%s\n%s", output, err)
			j.handleError(j.createError(formatted_error.Error(), cmd, status))

			return
		}
	}

}

var osRemoveAll = os.RemoveAll

func (j *Job) runInstallCmd() ([]byte, string, error) {
	j.nugetCommand = nuget
	installCmd, err := j.cmdFactory.MakeInstallCmd(j.nugetCommand, j.GetFile())

	if err != nil {
		command := ""
		if installCmd != nil {
			command = installCmd.String()
		}

		return nil, command, err
	}

	installCmdOutput, err := installCmd.Output()
	if err != nil {
		return installCmdOutput, installCmd.String(), j.GetExitError(err, "")
	}

	return installCmdOutput, installCmd.String(), nil
}

func (j *Job) cleanupTempCsproj() {
	// Cleanup of the temporary .csproj file (packages.config)
	tempFile := j.cmdFactory.GetTempoCsproj()
	if tempFile != "" {
		// remove the packages.config.csproj file
		err := osRemoveAll(tempFile)
		formatted_error := fmt.Errorf("failed to remove temporary .csproj file: %s", err)
		if err != nil {
			j.handleError(j.createError(formatted_error.Error(), "", "cleanup"))
		}
	}
}

func (j *Job) createError(error string, cmd string, status string) job.IError {
	cmdError := util.NewPMJobError(error)
	cmdError.SetCommand(cmd)
	cmdError.SetStatus(status)

	return cmdError
}

func (j *Job) handleError(cmdError job.IError) {
	expressions := []string{
		executableNotFoundErrRegex,
		versionNotFoundErrRegex,
		dependencyNotFoundErrRegex,
		unableToResolveErrRegex,
		noInternetErrRegex,
	}

	for _, expression := range expressions {
		regex := regexp.MustCompile(expression)
		matches := regex.FindAllStringSubmatch(cmdError.Error(), -1)

		if len(matches) > 0 {
			cmdError = j.addDocumentation(expression, matches, cmdError)
			j.Errors().Append(cmdError)

			return
		}
	}

	j.Errors().Append(cmdError)

	// Remove lock file
	dir := filepath.Dir(j.GetFile())
	os.Remove(filepath.Join(dir, packagesConfigLockfile))
}

func (j *Job) addDocumentation(expr string, matches [][]string, cmdError job.IError) job.IError {
	documentation := cmdError.Documentation()

	switch expr {
	case executableNotFoundErrRegex:
		documentation = j.GetExecutableNotFoundErrorDocumentation("Dotnet")
	case versionNotFoundErrRegex:
		documentation = j.getVersionNotFoundErrorDocumentation(matches)
	case dependencyNotFoundErrRegex:
		documentation = j.getDependencyNotFoundErrorDocumentation(matches)
	case unableToResolveErrRegex:
		documentation = j.getUnableToResolveErrorDocumentation(matches)
	case noInternetErrRegex:
		documentation = j.getNoInternetErrorDocumentation(matches)
	}

	cmdError.SetDocumentation(documentation)

	return cmdError
}

func (j *Job) getVersionNotFoundErrorDocumentation(matches [][]string) string {
	dependency := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		dependency = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Failed to find package",
			dependency + ".",
			"Please check that package versions are correct in the manifest file.",
		}, " ")
}

func (j *Job) getDependencyNotFoundErrorDocumentation(matches [][]string) string {
	dependency := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		dependency = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Failed to find package",
			"\"" + dependency + "\"",
			"that satisfies the requirements.",
			"Please check that dependencies are correct in the manifest file.",
			"\n" + util.InstallPrivateDependencyMessage,
		}, " ")
}

func (j *Job) getUnableToResolveErrorDocumentation(matches [][]string) string {
	dependency := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		dependency = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Couldn't resolve",
			"\"" + dependency + "\".",
			"Please check if it exists and NuGet sources are configured properly.",
		}, " ")
}

func (j *Job) getNoInternetErrorDocumentation(matches [][]string) string {
	registry := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		registry = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Registry",
			"\"" + registry + "\"",
			"is not available at the moment.",
			"There might be a trouble with your network connection,",
			"or this could be an authentication issue if this is a private registry.",
		}, " ")
}
