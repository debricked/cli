package maven

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	executableNotFoundErrRegex  = `executable file not found`
	lockFileExtension           = "maven.debricked.lock"
	nonParseablePomErrRegex     = "Non-parseable POM (.*)"
	networkUnreachableErrRegex  = "Failed to retrieve plugin descriptor"
	invalidVersionErrRegex      = "('[\\w\\.]+' for [\\w\\.:-]+ must not contain any of these characters .* but found .)"
	dependenciesResolveErrRegex = `Could not resolve dependencies for project\s+([\w\.-]+:[\w\.-]+:[\w\.-]+:[\w\.-]+)`
)

type Job struct {
	job.BaseJob
	cmdFactory ICmdFactory
	pomService IPomService
}

func NewJob(file string, cmdFactory ICmdFactory, pomService IPomService) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(file),
		cmdFactory: cmdFactory,
		pomService: pomService,
	}
}

func (j *Job) Run() {
	status := "parsing XML"
	j.SendStatus(status)

	file := j.GetFile()
	_, err := j.pomService.ParsePomModules(file)

	if err != nil {
		doc := err.Error()

		if doc == "EOF" {
			doc = "This file doesn't contain valid XML"
		}

		parsingError := util.NewPMJobError(err.Error())
		parsingError.SetStatus(status)
		parsingError.SetDocumentation(doc)

		j.Errors().Critical(parsingError)

		return
	}

	workingDirectory := filepath.Dir(filepath.Clean(file))
	cmd, err := j.cmdFactory.MakeDependencyTreeCmd(workingDirectory)
	if err != nil {
		j.handleError(util.NewPMJobError(err.Error()))

		return
	}

	status = "creating dependency graph"
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

func (j *Job) handleError(cmdError job.IError) {
	expressions := []string{
		executableNotFoundErrRegex,
		nonParseablePomErrRegex,
		networkUnreachableErrRegex,
		invalidVersionErrRegex,
		dependenciesResolveErrRegex,
	}

	for _, expression := range expressions {
		regex := regexp.MustCompile(expression)
		matches := regex.FindAllStringSubmatch(cmdError.Error(), -1)

		if len(matches) > 0 {
			cmdError = j.addDocumentation(expression, matches, cmdError)
			j.Errors().Critical(cmdError)

			return
		}
	}

	j.Errors().Critical(cmdError)
}

func (j *Job) addDocumentation(expr string, matches [][]string, cmdError job.IError) job.IError {
	documentation := cmdError.Documentation()

	switch expr {
	case executableNotFoundErrRegex:
		documentation = j.GetExecutableNotFoundErrorDocumentation("Maven")
	case nonParseablePomErrRegex:
		documentation = j.addNonParseablePomErrorDocumentation(matches)
	case networkUnreachableErrRegex:
		documentation = j.addNetworkUnreachableErrorDocumentation()
	case invalidVersionErrRegex:
		documentation = j.addInvalidVersionErrorDocumentation(matches)
	case dependenciesResolveErrRegex:
		documentation = j.addDependenciesResolveErrorDocumentation(matches)
	}

	cmdError.SetDocumentation(documentation)

	return cmdError
}

func (j *Job) addNonParseablePomErrorDocumentation(matches [][]string) string {
	message := "the POM file for errors"
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Failed to build Maven dependency tree.",
			"Your POM file is not valid.",
			"Please check",
			message,
		}, " ")
}

func (j *Job) addNetworkUnreachableErrorDocumentation() string {
	return strings.Join(
		[]string{
			"We weren't able to retrieve one or more plugin descriptor(s).",
			"Please check your Internet connection and try again.",
		}, " ")
}

func (j *Job) addInvalidVersionErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"There is an error in dependencies:",
			message,
		}, " ")
}

func (j *Job) addDependenciesResolveErrorDocumentation(matches [][]string) string {
	message := "An error occurred during dependencies resolve "
	if len(matches) > 0 && len(matches[0]) > 1 {
		message += strings.Join(
			[]string{
				"for: ",
				matches[0][1],
				"",
			}, "")
	}

	return strings.Join(
		[]string{
			message,
			"\nTry to run `mvn dependency:tree -e` to get more details.\n",
			util.InstallPrivateDependencyMessage,
		}, "")
}
