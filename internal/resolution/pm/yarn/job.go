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
	dependencyNotFoundErrRegex  = "error Error: (.*): Not found"
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
			cmdError := util.NewPMJobError(err.Error())
			cmdError.SetCommand(installCmd.String())
			cmdError.SetStatus(status)
			j.handleError(cmdError)

			return
		}

		_, err = installCmd.Output()

		if err != nil {
			cmdError := util.NewPMJobError(err.Error())
			cmdError.SetCommand(installCmd.String())
			cmdError.SetStatus(status)
			j.handleError(cmdError)

			return
		}
	}
}

func (j *Job) handleError(cmdErr job.IError) {
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
		matches := regex.FindAllStringSubmatch(cmdErr.Error(), -1)

		if len(matches) > 0 {
			cmdErr = j.addDocumentation(expression, matches, cmdErr)
			j.Errors().Append(cmdErr)

			return
		}
	}

	j.Errors().Append(cmdErr)
}

func (j *Job) addDocumentation(expr string, matches [][]string, cmdErr job.IError) job.IError {
	documentation := cmdErr.Documentation()

	switch {
	case expr == invalidJsonErrRegex:
		documentation = getInvalidJsonErrorDocumentation(matches, cmdErr)
	case expr == invalidSchemaErrRegex:
		documentation = getInvalidSchemaErrorDocumentation(matches, cmdErr)
	case expr == invalidArgumentErrRegex:
		documentation = getInvalidArgumentErrorDocumentation(matches, cmdErr)
	case expr == versionNotFoundErrRegex:
		documentation = getVersionNotFoundErrorDocumentation(matches, cmdErr)
	case expr == dependencyNotFoundErrRegex:
		documentation = getDependencyNotFoundErrorDocumentation(matches, cmdErr)
	case expr == registryUnavailableErrRegex:
		documentation = getRegistryUnavailableErrorDocumentation(matches, cmdErr)
	case expr == permissionDeniedErrRegex:
		documentation = getPermissionDeniedErrorDocumentation(matches, cmdErr)
	}

	cmdErr.SetDocumentation(documentation)

	return cmdErr
}

func getInvalidJsonErrorDocumentation(matches [][]string, err job.IError) string {
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

func getInvalidSchemaErrorDocumentation(matches [][]string, err job.IError) string {
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

func getInvalidArgumentErrorDocumentation(matches [][]string, err job.IError) string {
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

func getDependencyNotFoundErrorDocumentation(matches [][]string, err job.IError) string {
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
		}, " ")
}

func getVersionNotFoundErrorDocumentation(matches [][]string, err job.IError) string {
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

func getRegistryUnavailableErrorDocumentation(matches [][]string, err job.IError) string {
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

func getPermissionDeniedErrorDocumentation(matches [][]string, err job.IError) string {
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
