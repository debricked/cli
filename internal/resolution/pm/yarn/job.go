package yarn

import (
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	yarn                        = "yarn"
	invalidJsonErrRegex         = "error SyntaxError.*package.json: (.*)"
	invalidSchemaErrRegex       = "error package.json: (.*)"
	invalidArgumentErrRegex     = "error TypeError \\[\\w+\\]: (.*)"
	versionNotFoundErrRegex     = "error (Couldn\\'t find any versions for .*)"
	dependencyNotFoundErrRegex  = `error.*? "?(https?://[^"\s:]+)?: Not found`
	registryUnavailableErrRegex = "error Error: getaddrinfo ENOTFOUND ([\\w\\.]+)"
	permissionDeniedErrRegex    = "Error: (.*): Request failed \"404 Not Found\""
)

type Job struct {
	job.BaseJob
	install     bool
	yarnCommand string
	cmdFactory  ICmdFactory
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
		j.yarnCommand = yarn

		installCmd, err := j.cmdFactory.MakeInstallCmd(j.yarnCommand, j.GetFile())

		if err != nil {
			j.handleError(j.createError(err.Error(), installCmd.String(), status))

			return
		}

		if output, err := installCmd.Output(); err != nil {
			error := strings.Join([]string{string(output), j.GetExitError(err).Error()}, "")
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
		invalidJsonErrRegex,
		invalidSchemaErrRegex,
		invalidArgumentErrRegex,
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

	switch {
	case expr == invalidJsonErrRegex:
		documentation = getInvalidJsonErrorDocumentation(matches)
	case expr == invalidSchemaErrRegex:
		documentation = getInvalidSchemaErrorDocumentation(matches)
	case expr == invalidArgumentErrRegex:
		documentation = getInvalidArgumentErrorDocumentation(matches)
	case expr == versionNotFoundErrRegex:
		documentation = getVersionNotFoundErrorDocumentation(matches)
	case expr == dependencyNotFoundErrRegex:
		documentation = getDependencyNotFoundErrorDocumentation(matches)
	case expr == registryUnavailableErrRegex:
		documentation = getRegistryUnavailableErrorDocumentation(matches)
	case expr == permissionDeniedErrRegex:
		documentation = getPermissionDeniedErrorDocumentation(matches)
	}

	cmdError.SetDocumentation(documentation)

	return cmdError
}

func getInvalidJsonErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Your package.json file contains invalid JSON:",
			message + ".",
		}, " ")
}

func getInvalidSchemaErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Your package.json file is not valid:",
			message + ".",
			"Please make sure it follows the schema.",
		}, " ")
}

func getInvalidArgumentErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			message + ".",
			"Please make sure that your package.json file doesn't contain errors.",
		}, " ")
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
			"that satisfies the requirement from yarn dependencies.",
			"Please check that dependencies are correct in your package.json file.",
			"\n" + util.InstallPrivateDependencyMessage,
		}, " ")
}

func getVersionNotFoundErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			message + ".",
			"Please check that dependencies are correct in your package.json file.",
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
	dependency := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		dependency = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Failed to find a package that satisfies requirements for yarn dependencies:",
			dependency + ".",
			"This could mean that the package or version does not exist or is private.\n",
			util.InstallPrivateDependencyMessage,
		}, " ")
}
