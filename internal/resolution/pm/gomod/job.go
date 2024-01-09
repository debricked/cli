package gomod

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/debricked/cli/internal/resolution/pm/writer"
)

const (
	fileName                   = "gomod.debricked.lock"
	versionNotFoundErrRegex    = `require ([^"'\s:]+): version "[^"'\s:]+" invalid: ([^"'\n:]+)`
	revisionNotFoundErrRegex   = `([^"'\s\n:]+): reading [^"'\n:]+ at revision [^"'\n:]+: unknown revision ([^"'\n:]+)`
	dependencyNotFoundErrRegex = `go: ([^"'\s:]+): .*\n.*fatal: could not read Username`
	repositoryNotFoundErrRegex = `go: ([^"'\s:]+): .*\n.*remote: Repository not found`
	noPackageErrRegex          = `([^"'\s:]+): .*, but does not contain package`
	unableToResolveErrRegex    = `go: module ([^"'\s:]+): .*\n.*Permission denied`
	noInternetErrRegex         = `dial tcp: lookup ([^"'\s:]+) .+: server misbehaving`
)

type Job struct {
	job.BaseJob
	cmdFactory ICmdFactory
	fileWriter writer.IFileWriter
}

func NewJob(
	file string,
	cmdFactory ICmdFactory,
	fileWriter writer.IFileWriter,
) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(file),
		cmdFactory: cmdFactory,
		fileWriter: fileWriter,
	}
}

func (j *Job) Run() {
	status := "creating dependency graph"
	j.SendStatus(status)

	workingDirectory := filepath.Dir(filepath.Clean(j.GetFile()))

	graphCmdOutput, cmd, err := j.runGraphCmd(workingDirectory)
	if err != nil {
		j.handleError(j.createError(err.Error(), cmd, status))

		return
	}

	status = "creating dependency version list"
	j.SendStatus(status)
	listCmdOutput, cmd, err := j.runListCmd(workingDirectory)
	if err != nil {
		j.handleError(j.createError(err.Error(), cmd, status))

		return
	}

	status = "creating lock file"
	j.SendStatus(status)
	lockFile, err := j.fileWriter.Create(util.MakePathFromManifestFile(j.GetFile(), fileName))
	if err != nil {
		j.handleError(j.createError(err.Error(), cmd, status))

		return
	}
	defer util.CloseFile(j, j.fileWriter, lockFile)

	var fileContents []byte
	fileContents = append(fileContents, graphCmdOutput...)
	fileContents = append(fileContents, []byte("\n")...)
	fileContents = append(fileContents, listCmdOutput...)

	err = j.fileWriter.Write(lockFile, fileContents)
	if err != nil {
		j.handleError(j.createError(err.Error(), cmd, status))
	}
}

func (j *Job) runGraphCmd(workingDirectory string) ([]byte, string, error) {
	graphCmd, err := j.cmdFactory.MakeGraphCmd(workingDirectory)
	if err != nil {
		return nil, graphCmd.String(), err
	}

	graphCmdOutput, err := graphCmd.Output()
	if err != nil {
		return nil, graphCmd.String(), j.GetExitError(err, "")
	}

	return graphCmdOutput, graphCmd.String(), nil
}

func (j *Job) runListCmd(workingDirectory string) ([]byte, string, error) {
	listCmd, err := j.cmdFactory.MakeListCmd(workingDirectory)
	if err != nil {
		return nil, listCmd.String(), err
	}

	listCmdOutput, err := listCmd.Output()
	if err != nil {
		return nil, listCmd.String(), j.GetExitError(err, "")
	}

	return listCmdOutput, listCmd.String(), nil
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
		revisionNotFoundErrRegex,
		dependencyNotFoundErrRegex,
		repositoryNotFoundErrRegex,
		noPackageErrRegex,
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
}

func (j *Job) addDocumentation(expr string, matches [][]string, cmdError job.IError) job.IError {
	documentation := cmdError.Documentation()

	switch expr {
	case versionNotFoundErrRegex:
		documentation = getVersionNotFoundErrorDocumentation(matches)
	case revisionNotFoundErrRegex:
		documentation = getRevisionNotFoundErrorDocumentation(matches)
	case dependencyNotFoundErrRegex:
		documentation = getDependencyNotFoundErrorDocumentation(matches)
	case repositoryNotFoundErrRegex:
		documentation = getDependencyNotFoundErrorDocumentation(matches)
	case noPackageErrRegex:
		documentation = getNoPackageErrorDocumentation(matches)
	case unableToResolveErrRegex:
		documentation = getDependencyNotFoundErrorDocumentation(matches)
	case noInternetErrRegex:
		documentation = getNoInternetErrorDocumentation(matches)
	}

	cmdError.SetDocumentation(documentation)

	return cmdError
}

func getVersionNotFoundErrorDocumentation(matches [][]string) string {
	dependency := ""
	recommendation := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		dependency = matches[0][1]
		recommendation = matches[0][2]
	}

	return strings.Join(
		[]string{
			"Failed to find package",
			"\"" + dependency + "\".",
			"Please check that package versions are correct in the manifest file.",
			"It " + recommendation + ".",
		}, " ")
}

func getRevisionNotFoundErrorDocumentation(matches [][]string) string {
	dependency := ""
	revision := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		dependency = matches[0][1]
		revision = matches[0][2]
	}

	return strings.Join(
		[]string{
			"Failed to find package",
			"\"" + dependency + "\" with revision " + revision + ".",
			"Please check that package version is correct in the manifest file.",
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
			"that satisfies the requirements.",
			"Please check that dependencies are correct in the manifest file.",
			"\n" + util.InstallPrivateDependencyMessage,
		}, " ")
}

func getNoPackageErrorDocumentation(matches [][]string) string {
	repository := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		repository = matches[0][1]
	}

	return strings.Join(
		[]string{
			"We weren't able to find a package in provided repository",
			"\"" + repository + "\".",
			"Please check that repository address is spelled correct and it actually contains a Go package.",
		}, " ")
}

func getNoInternetErrorDocumentation(matches [][]string) string {
	registry := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		registry = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Registry",
			"\"" + registry + "\"",
			"is not available at the moment.",
			"There might be a trouble with your network connection.",
		}, " ")
}
