package pub

import (
	"os"
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	executableNotFoundErrRegex = `executable file not found`
	depsFileName               = "pubspec.deps.json"
)

// Job runs `dart pub get` and `dart pub deps --json` for a given pubspec.yaml.
// It produces pubspec.lock and pubspec.deps.json, where the latter contains
// explicit parent-child relationships for full transitive tree reconstruction.
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
	status := "generating pubspec.lock"
	j.SendStatus(status)

	lockCmd, err := j.cmdFactory.MakeLockCmd(j.GetFile())
	if err != nil {
		j.handleError(j.createError(err.Error(), "", status))

		return
	}

	if output, err := lockCmd.Output(); err != nil {
		exitErr := j.GetExitError(err, string(output))
		errorMessage := strings.Join([]string{string(output), exitErr.Error()}, "")
		j.handleError(j.createError(errorMessage, lockCmd.String(), status))

		return
	}

	status = "generating pubspec.deps.json"
	j.SendStatus(status)

	depsCmd, err := j.cmdFactory.MakeDepsCmd(j.GetFile())
	if err != nil {
		j.handleError(j.createError(err.Error(), "", status))

		return
	}

	depsOutput, err := depsCmd.Output()
	if err != nil {
		exitErr := j.GetExitError(err, string(depsOutput))
		errorMessage := strings.Join([]string{string(depsOutput), exitErr.Error()}, "")
		j.handleError(j.createError(errorMessage, depsCmd.String(), status))

		return
	}

	status = "writing pubspec.deps.json"
	j.SendStatus(status)

	err = os.WriteFile(util.MakePathFromManifestFile(j.GetFile(), depsFileName), depsOutput, 0600)
	if err != nil {
		j.handleError(j.createError(err.Error(), "", status))

		return
	}
}

func (j *Job) createError(errorStr string, cmd string, status string) job.IError {
	cmdError := util.NewPMJobError(errorStr)
	cmdError.SetCommand(cmd)
	cmdError.SetStatus(status)

	return cmdError
}

func (j *Job) handleError(cmdError job.IError) {
	expressions := []string{
		executableNotFoundErrRegex,
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

func (j *Job) addDocumentation(expr string, _ [][]string, cmdError job.IError) job.IError {
	documentation := cmdError.Documentation()

	switch expr {
	case executableNotFoundErrRegex:
		documentation = j.GetExecutableNotFoundErrorDocumentation("Dart")
	}

	cmdError.SetDocumentation(documentation)

	return cmdError
}
