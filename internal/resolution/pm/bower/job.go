package bower

import (
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/debricked/cli/internal/resolution/pm/writer"
)

const (
	bower                       = "bower"
	fileName                    = "bower.debricked.lock"
	bowerNotFoundErrRegex       = `executable file not found`
	versionNotFoundErrRegex     = `([^"\s:]+)\s+ENORESTARGET No tag found`
	dependencyNotFoundErrRegex  = `ENOTFOUND Package ([^"\s:]+) not found`
	registryUnavailableErrRegex = `getaddrinfo EAI_AGAIN ([\w\/\-\.]+)`
	permissionDeniedErrRegex    = `EACCES: permission denied, rmdir '([\w\/-]+)'`
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
	status := "installing dependencies"
	j.SendStatus(status)

	cmd, err := j.runInstallCmd(j.GetFile())
	if err != nil {
		j.handleError(j.createError(err.Error(), cmd, status))

		return
	}

	status = "creating dependency tree"
	j.SendStatus(status)
	listCmdOutput, cmd, err := j.runListCmd(j.GetFile())
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

	err = j.fileWriter.Write(lockFile, listCmdOutput)
	if err != nil {
		j.handleError(j.createError(err.Error(), cmd, status))
	}
}

func (j *Job) runInstallCmd(file string) (string, error) {
	installCmd, err := j.cmdFactory.MakeInstallCmd(bower, file)
	if err != nil {
		return installCmd.String(), err
	}

	_, err = installCmd.Output()
	if err != nil {
		return installCmd.String(), j.GetExitError(err, "")
	}

	return installCmd.String(), nil
}

func (j *Job) runListCmd(file string) ([]byte, string, error) {
	listCmd, err := j.cmdFactory.MakeListCmd(bower, file)
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
		bowerNotFoundErrRegex,
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
	case bowerNotFoundErrRegex:
		documentation = getBowerNotFoundErrorDocumentation()
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

func getBowerNotFoundErrorDocumentation() string {
	return strings.Join(
		[]string{
			"Bower wasn't found.",
			"Please check if it is installed and accessible by the CLI.",
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
			"Please check that dependencies are correct in your bower.json file.",
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
			"that satisfies the requirement from bower.json file.",
			"In most cases you or one of your dependencies are requesting a package version that doesn't exist.",
			"Please check that package versions are correct.",
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
