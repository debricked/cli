package npm

import (
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	npm                         = "npm"
	versionNotFoundErrRegex     = `notarget [\w\s]+ ([^"\s:]+).`
	dependencyNotFoundErrRegex  = `404\s+'([^"\s:]+)'`
	registryUnavailableErrRegex = `EAI_AGAIN ([\w\.]+)`
	permissionDeniedErrRegex    = `Error: EACCES, open '([^"\s:]+)'`
)

type Job struct {
	job.BaseJob
	install    bool
	npmCommand string
	cmdFactory ICmdFactory
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
		j.npmCommand = npm

		installCmd, err := j.cmdFactory.MakeInstallCmd(j.npmCommand, j.GetFile())

		if err != nil {
			j.handleError(j.createError(err.Error(), installCmd.String(), status))

			return
		}

		if output, err := installCmd.Output(); err != nil {
			error := strings.Join([]string{string(output), j.GetExitError(err, "").Error()}, "")
			j.handleError(j.createError(error, installCmd.String(), status))

			return
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
		versionNotFoundErrRegex,
		dependencyNotFoundErrRegex,
		registryUnavailableErrRegex,
		permissionDeniedErrRegex,
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
}

func (j *Job) addDocumentation(expr string, matches [][]string, cmdError job.IError) job.IError {
	documentation := cmdError.Documentation()

	switch expr {
	case versionNotFoundErrRegex:
		documentation = getVersionNotFoundErrorDocumentation(matches)
	case dependencyNotFoundErrRegex:
		documentation = getDependencyNotFoundErrorDocumentation(matches)
	case registryUnavailableErrRegex:
		documentation = getRegistryUnavailableErrorDocumentation(matches)
	case permissionDeniedErrRegex:
		documentation = getPermissionDeniedErrorDocumentation(matches)
	}

	cmdError.SetDocumentation(documentation)

	return cmdError
}

func getDependencyNotFoundErrorDocumentation(matches [][]string) string {
	dependency := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		dependency = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Failed to find package",
			"\"" + dependency + "\"",
			"that satisfies the requirement from dependencies.",
			"Please check that dependencies are correct in your package.json file.",
			"\n" + util.InstallPrivateDependencyMessage,
		}, " ")
}

func getVersionNotFoundErrorDocumentation(matches [][]string) string {
	dependency := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		dependency = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Failed to find package",
			"\"" + dependency + "\"",
			"that satisfies the requirement from package.json file.",
			"In most cases you or one of your dependencies are requesting a package version that doesn't exist.",
			"Please check that package versions are correct in your package.json file.",
		}, " ")
}

func getRegistryUnavailableErrorDocumentation(matches [][]string) string {
	registry := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		registry = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Package registry",
			"\"" + registry + "\"",
			"is not available at the moment.",
			"There might be a trouble with your network connection.",
		}, " ")
}

func getPermissionDeniedErrorDocumentation(matches [][]string) string {
	path := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		path = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Couldn't get access to",
			"\"" + path + "\".",
			"Please check permissions or try running this command again as root/Administrator.",
		}, " ")
}
