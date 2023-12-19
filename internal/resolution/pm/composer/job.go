package composer

import (
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	composer                    = "composer"
	composerMissingExtension    = "Composer requires it to run"
	invalidRequirement          = `require\.([^ ]*) is invalid, it should have a vendor name`
	noNetworkRegex              = `The following exception probably indicates you( are offline or)? have misconfigured DNS resolver\(s\)`
	invalidVersionErrRegex      = `requires\s+([^/]+/[^,]+),\s+found.*but it does not match the constraint\.`
	dependenciesResolveErrRegex = `requires\s+([^/]+/[^,]+),\s+it\s+could\s+not\s+be\s+found\s+in\s+any\s+version`
)

type Job struct {
	job.BaseJob
	install         bool
	composerCommand string
	cmdFactory      ICmdFactory
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
		_, err := j.runInstallCmd()
		if err != nil {
			cmdErr := util.NewPMJobError(err.Error())
			j.handleError(cmdErr)

			return
		}
	}

}

func (j *Job) runInstallCmd() ([]byte, error) {
	j.composerCommand = composer
	installCmd, err := j.cmdFactory.MakeInstallCmd(j.composerCommand, j.GetFile())
	if err != nil {
		return nil, err
	}

	installCmdOutput, err := installCmd.Output()
	if err != nil {
		return nil, j.GetExitError(err, string(installCmdOutput))
	}

	return installCmdOutput, nil
}

func (j *Job) handleError(cmdErr job.IError) {
	expressions := []string{
		composerMissingExtension,
		invalidRequirement,
		noNetworkRegex,
		invalidVersionErrRegex,
		dependenciesResolveErrRegex,
	}

	for _, expression := range expressions {
		regex := regexp.MustCompile(expression)

		if regex.MatchString(cmdErr.Error()) {
			cmdErr = j.addDocumentation(expression, regex, cmdErr)
			j.Errors().Critical(cmdErr)

			return
		}
	}

	j.Errors().Critical(cmdErr)
}

func (j *Job) addDocumentation(expr string, regex *regexp.Regexp, cmdErr job.IError) job.IError {
	switch expr {
	case composerMissingExtension:
		cmdErr = j.addComposerMissingRequirementsErrorDocumentation(regex, cmdErr)
	case invalidRequirement:
		cmdErr = j.addInvalidRequirementErrorDocumentation(regex, cmdErr)
	case noNetworkRegex:
		cmdErr = j.addNetworkUnreachableErrorDocumentation(cmdErr)
	case invalidVersionErrRegex:
		cmdErr = j.addInvalidVersionErrorDocumentation(regex, cmdErr)
	case dependenciesResolveErrRegex:
		cmdErr = j.addDependenciesResolveErrorDocumentation(regex, cmdErr)
	}

	return cmdErr
}

func (j *Job) addComposerMissingRequirementsErrorDocumentation(regex *regexp.Regexp, cmdErr job.IError) job.IError {
	cmdErr.SetDocumentation(
		strings.Join(
			[]string{
				"Failed to build Composer dependency tree.",
				"Your runtime environment is missing one or more Composer requirements.",
				"Check error message below for more details:\n\n",
				cmdErr.Error(),
			}, " "),
	)

	return cmdErr
}

func (j *Job) addInvalidRequirementErrorDocumentation(regex *regexp.Regexp, cmdErr job.IError) job.IError {
	matches := regex.FindAllStringSubmatch(cmdErr.Error(), 1)
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	cmdErr.SetDocumentation(
		strings.Join(
			[]string{
				"Couldn't resolve dependency",
				message,
				", please make sure it is spelt correctly:\n",
			}, " "),
	)

	return cmdErr
}

func (j *Job) addNetworkUnreachableErrorDocumentation(cmdErr job.IError) job.IError {
	cmdErr.SetDocumentation(
		strings.Join(
			[]string{
				"We weren't able to retrieve one or more dependencies.",
				"Please check your Internet connection and try again.",
			}, " "),
	)

	return cmdErr
}

func (j *Job) addInvalidVersionErrorDocumentation(regex *regexp.Regexp, cmdErr job.IError) job.IError {
	matches := regex.FindAllStringSubmatch(cmdErr.Error(), 1)
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	cmdErr.SetDocumentation(
		strings.Join(
			[]string{
				"Couldn't resolve version",
				message,
				", please make sure it exists:\n",
			}, " "),
	)

	return cmdErr
}

func (j *Job) addDependenciesResolveErrorDocumentation(regex *regexp.Regexp, cmdErr job.IError) job.IError {
	matches := regex.FindAllStringSubmatch(cmdErr.Error(), 1)
	message := "An error occurred during dependencies resolve "
	if len(matches) > 0 && len(matches[0]) > 1 {
		message += strings.Join(
			[]string{
				"for: ",
				matches[0][1],
				"",
			}, "")
	}

	cmdErr.SetDocumentation(
		strings.Join(
			[]string{
				message,
				"\n\n",
				util.InstallPrivateDependencyMessage,
				"\n\n",
			}, ""),
	)

	return cmdErr
}
