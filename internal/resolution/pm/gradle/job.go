package gradle

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/debricked/cli/internal/resolution/pm/writer"
)

const (
	bugErrRegex                = "BUG! (.*)"
	executableNotFoundErrRegex = `executable file not found`
	notRootDirErrRegex         = "Error: (Could not find or load main class .*)"
	unrelatedBuildErrRegex     = "(Project directory '.*' is not part of the build defined by settings file '.*')"
	unknownPropertyErrRegex    = "(Could not get unknown property .*)"
)

type Job struct {
	job.BaseJob
	dir              string
	gradlew          string
	groovyInitScript string
	cmdFactory       ICmdFactory
	fileWriter       writer.IFileWriter
}

func NewJob(
	file string,
	dir string,
	gradlew string,
	groovyInitScript string,
	cmdFactory ICmdFactory,
	fileWriter writer.IFileWriter,
) *Job {

	return &Job{
		BaseJob:          job.NewBaseJob(file),
		dir:              dir,
		gradlew:          gradlew,
		groovyInitScript: groovyInitScript,
		cmdFactory:       cmdFactory,
		fileWriter:       fileWriter,
	}
}

func (j *Job) Run() {
	workingDirectory := filepath.Clean(j.GetDir())
	dependenciesCmd, err := j.cmdFactory.MakeDependenciesGraphCmd(workingDirectory, j.gradlew, j.groovyInitScript)
	var permissionErr error

	if err != nil {
		if strings.HasSuffix(err.Error(), "gradlew\": permission denied") {
			permissionErr = fmt.Errorf("Permission to execute gradlew is not granted, fallback to PATHs gradle installation will be used.\nFull error: %s", err.Error())

			dependenciesCmd, err = j.cmdFactory.MakeDependenciesGraphCmd(workingDirectory, "gradle", j.groovyInitScript)
		}
	}

	if err != nil {
		if permissionErr != nil {
			j.handleError(util.NewPMJobError(permissionErr.Error()))
		}
		cmdErr := util.NewPMJobError(err.Error())
		cmdErr.SetCommand(strings.Trim(dependenciesCmd.String(), " "))
		j.handleError(cmdErr)

		return
	}

	status := "creating dependency graph"
	j.SendStatus(status)
	_, err = dependenciesCmd.Output()

	if permissionErr != nil {
		cmdErr := util.NewPMJobError(permissionErr.Error())
		cmdErr.SetIsCritical(false)
		j.handleError(cmdErr)
	}

	if err != nil {
		cmdErr := util.NewPMJobError(j.GetExitError(err, "").Error())
		cmdErr.SetCommand(dependenciesCmd.String())
		j.handleError(cmdErr)

		return
	}
}

func (j *Job) GetDir() string {
	return j.dir
}

func (j *Job) handleError(cmdError job.IError) {
	expressions := []string{
		bugErrRegex,
		executableNotFoundErrRegex,
		notRootDirErrRegex,
		unrelatedBuildErrRegex,
		unknownPropertyErrRegex,
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
	case executableNotFoundErrRegex:
		documentation = j.GetExecutableNotFoundErrorDocumentation("Gradle")
	case bugErrRegex:
		documentation = j.addBugErrorDocumentation(matches)
	case notRootDirErrRegex:
		documentation = j.addNotRootDirErrorDocumentation(matches)
	case unrelatedBuildErrRegex:
		documentation = j.addUnrelatedBuildErrorDocumentation(matches)
	case unknownPropertyErrRegex:
		documentation = j.addUnknownPropertyErrorDocumentation(matches)
	}

	cmdError.SetDocumentation(documentation)

	return cmdError
}

func (j *Job) addBugErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Failed to build Gradle dependency tree. ",
			"The process has failed with following error: ",
			message,
			". ",
			"Try running the command below with --stacktrace flag to get a stacktrace. ",
			"Replace --stacktrace with --info or --debug option to get more log output. ",
			"Or with --scan to get full insights.",
		}, "")
}

func (j *Job) addNotRootDirErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Failed to build Gradle dependency tree.",
			"The process has failed with following error: " + message + ".", //nolint:all
			"You are probably not running the command from the root directory.",
		}, " ")
}

func (j *Job) addUnrelatedBuildErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Failed to build Gradle dependency tree. ",
			"The process has failed with following error: ",
			message,
			". ",
			"This error might be caused by inclusion of test folders into resolve process. ",
			"Try running resolve command with -e flag. ",
			"For example, `debricked resolve -e \"**/test*/**\"` will exclude all folders that start from 'test' from resolution process. ",
			"Or if this is an unrelated build, it must have its own settings file.",
		}, "")
}

func (j *Job) addUnknownPropertyErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Failed to build Gradle dependency tree. ",
			"The process has failed with following error: ",
			message,
			". ",
			"Please check your settings.gradle file for errors.",
		}, "")
}
