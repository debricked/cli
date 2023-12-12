package maven

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	lockFileExtension           = "maven.debricked.lock"
	nonParseablePomErrRegex     = "Non-parseable POM (.*)"
	networkUnreachableErrRegex  = "Failed to retrieve plugin descriptor"
	invalidVersionErrRegex      = "('[\\w\\.]+' for [\\w\\.:-]+ must not contain any of these characters .* but found .)"
	dependenciesResolveErrRegex = `Could not resolve dependencies for project`
)

type Job struct {
	job.BaseJob
	cmdFactory ICmdFactory
}

func NewJob(file string, cmdFactory ICmdFactory) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(file),
		cmdFactory: cmdFactory,
	}
}

func (j *Job) Run() {
	workingDirectory := filepath.Dir(filepath.Clean(j.GetFile()))
	cmd, err := j.cmdFactory.MakeDependencyTreeCmd(workingDirectory)
	if err != nil {
		j.handleError(util.NewPMJobError(err.Error()))

		return
	}

	status := "creating dependency graph"
	j.SendStatus(status)
	var output []byte
	output, err = cmd.Output()
	if err != nil {
		errContent := err.Error()
		if output != nil {
			errContent = string(output)
		}

		cmdErr := util.NewPMJobError(errContent)
		cmdErr.SetStatus(status)

		j.handleError(cmdErr)
	}
}

func (j *Job) handleError(cmdErr job.IError) {
	expressions := []string{
		nonParseablePomErrRegex,
		networkUnreachableErrRegex,
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
	switch {
	case expr == nonParseablePomErrRegex:
		cmdErr = j.addNonParseablePomErrorDocumentation(regex, cmdErr)
	case expr == networkUnreachableErrRegex:
		cmdErr = j.addNetworkUnreachableErrorDocumentation(cmdErr)
	case expr == invalidVersionErrRegex:
		cmdErr = j.addInvalidVersionErrorDocumentation(regex, cmdErr)
	case expr == dependenciesResolveErrRegex:
		cmdErr = j.addDependenciesResolveErrorDocumentation(regex, cmdErr)
	}

	return cmdErr
}

func (j *Job) addNonParseablePomErrorDocumentation(regex *regexp.Regexp, cmdErr job.IError) job.IError {
	matches := regex.FindAllStringSubmatch(cmdErr.Error(), 1)
	message := "the POM file for errors"
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	cmdErr.SetDocumentation(
		strings.Join(
			[]string{
				"Failed to build Maven dependency tree.",
				"Your POM file is not valid.",
				"Please check",
				message,
			}, " "),
	)

	return cmdErr
}

func (j *Job) addNetworkUnreachableErrorDocumentation(cmdErr job.IError) job.IError {
	cmdErr.SetDocumentation(
		strings.Join(
			[]string{
				"We weren't able to retrieve one or more plugin descriptor(s).",
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
				"There is an error in dependencies:",
				message,
			}, " "),
	)

	return cmdErr
}

func (j *Job) addDependenciesResolveErrorDocumentation(regex *regexp.Regexp, cmdErr job.IError) job.IError {
	matches := regex.FindAllStringSubmatch(cmdErr.Error(), 1)
	message := "An error occurred during dependencies resolve"
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	cmdErr.SetDocumentation(
		strings.Join(
			[]string{
				message,
				"\nTry to run `mvn dependency:tree -e` to get more details.",
				util.InstallPrivateDependencyMessage,
			}, " "),
	)

	return cmdErr
}
