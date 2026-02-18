package pnpm

import (
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	pnpm                       = "pnpm"
	executableNotFoundErrRegex = `executable file not found`
)

type Job struct {
	job.BaseJob
	install     bool
	pnpmCommand string
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
		j.pnpmCommand = pnpm

		installCmd, err := j.cmdFactory.MakeInstallCmd(j.pnpmCommand, j.GetFile())

		if err != nil {
			j.handleError(j.createError(err.Error(), installCmd.String(), status))

			return
		}

		if output, err := installCmd.Output(); err != nil {
			joined := strings.Join([]string{string(output), j.GetExitError(err, "").Error()}, "")
			j.handleError(j.createError(joined, installCmd.String(), status))

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

func (j *Job) addDocumentation(expr string, matches [][]string, cmdError job.IError) job.IError {
	documentation := cmdError.Documentation()

	switch expr {
	case executableNotFoundErrRegex:
		documentation = j.GetExecutableNotFoundErrorDocumentation("PNPM")
	}

	cmdError.SetDocumentation(documentation)

	return cmdError
}
