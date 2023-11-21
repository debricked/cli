package pip

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/debricked/cli/internal/resolution/pm/writer"
	internalOs "github.com/debricked/cli/internal/runtime/os"
)

const (
	lockFileExtension = ".pip.debricked.lock"
	pip               = "pip"
	lockFileDelimiter = "***"
)

type Job struct {
	job.BaseJob
	install    bool
	venvPath   string
	pipCommand string
	cmdFactory ICmdFactory
	fileWriter writer.IFileWriter
	pipCleaner IPipCleaner
}

func NewJob(
	file string,
	install bool,
	cmdFactory ICmdFactory,
	fileWriter writer.IFileWriter,
	pipCleaner IPipCleaner,
) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(file),
		install:    install,
		cmdFactory: cmdFactory,
		fileWriter: fileWriter,
		pipCleaner: pipCleaner,
	}
}

type IPipCleaner interface {
	RemoveAll(path string) error
}

type pipCleaner struct{}

func (p pipCleaner) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (j *Job) Install() bool {
	return j.install
}

func (j *Job) Run() {
	if j.install {
		defer func() {
			if j.venvPath == "" {
				return
			}
			status := "removing venv"
			j.SendStatus(status)
			err := j.pipCleaner.RemoveAll(j.venvPath)
			if err != nil {
				cmdErr := util.NewPMJobError(err.Error())
				cmdErr.SetDocumentation("Error when trying to remove previous virtual environment")
				cmdErr.SetStatus(status)
				j.Errors().Critical(cmdErr)
			}
		}()
		status := "creating venv"
		j.SendStatus(status)
		_, cmdErr := j.runCreateVenvCmd()
		if cmdErr != nil {
			cmdErr.SetDocumentation("Error when trying to create python virtual environment")
			cmdErr.SetStatus(status)
			j.Errors().Critical(cmdErr)

			return
		}
		status = "installing dependencies"
		j.SendStatus(status)
		_, cmdErr = j.runInstallCmd()
		if cmdErr != nil {
			cmdErr.SetStatus(status)
			j.handleInstallError(cmdErr)

			return
		}
	}

	err := j.writeLockContent()
	if err != nil {
		j.Errors().Critical(err)

		return
	}

}

func (j *Job) handleInstallError(cmdErr job.IError) {
	var buildError = regexp.MustCompile("setup.py[ install for]*(?P<dependency>[^ ]*) did not run successfully.")
	var credentialError = regexp.MustCompile("WARNING: 401 Error, Credentials not correct for")

	switch {
	case buildError.MatchString(cmdErr.Error()):
		matches := buildError.FindAllStringSubmatch(cmdErr.Error(), 10)
		dependencyName := ""
		if len(matches) > 0 {
			if len(matches[len(matches)-1]) > 1 {
				dependencyName = "\"" + matches[len(matches)-1][1] + "\""
			}
		}
		cmdErr.SetDocumentation(
			strings.Join(
				[]string{
					"Failed to build python dependency ",
					dependencyName,
					" with setup.py. This probably means the " +
						"project was not set up correctly and " +
						"could mean that an OS package is missing.",
				}, ""),
		)
	case credentialError.MatchString(cmdErr.Error()):
		authErrDependencyNamePattern := regexp.MustCompile("No matching distribution found for ([^ ]*)\n")
		dependencyNameMatch := authErrDependencyNamePattern.FindStringSubmatch(cmdErr.Error())
		dependencyName := ""
		if len(dependencyNameMatch) > 1 {
			dependencyName = "\"" + dependencyNameMatch[len(dependencyNameMatch)-1] + "\""
		}
		cmdErr.SetDocumentation(
			strings.Join(
				[]string{
					"Failed to install python dependency ",
					dependencyName,
					" due to authorization. This could mean it is",
					" a private dependency that the debricked CLI",
					" does not have access to.",
				}, ""),
		)
	}
	j.Errors().Critical(cmdErr)
}

func (j *Job) writeLockContent() job.IError {
	status := "generating lock file"
	j.SendStatus(status)
	catCmdOutput, cmdErr := j.runCatCmd()
	if cmdErr != nil {
		cmdErr.SetStatus(status)

		return cmdErr
	}

	listCmdOutput, cmdErr := j.runListCmd()
	if cmdErr != nil {
		cmdErr.SetStatus(status)

		return cmdErr
	}

	installedPackages := j.parsePipList(string(listCmdOutput))
	ShowCmdOutput, cmdErr := j.runShowCmd(installedPackages)
	if cmdErr != nil {
		cmdErr.SetStatus(status)

		return cmdErr
	}

	lockFileName := fmt.Sprintf("%s%s", filepath.Base(j.GetFile()), lockFileExtension)
	lockFile, err := j.fileWriter.Create(util.MakePathFromManifestFile(j.GetFile(), lockFileName))
	if err != nil {
		cmdErr = util.NewPMJobError(err.Error())
		cmdErr.SetStatus(status)

		return cmdErr
	}
	defer closeFile(j, lockFile)

	var fileContents []string
	fileContents = append(fileContents, string(catCmdOutput))
	fileContents = append(fileContents, lockFileDelimiter)
	fileContents = append(fileContents, string(listCmdOutput))
	fileContents = append(fileContents, lockFileDelimiter)
	fileContents = append(fileContents, string(ShowCmdOutput))
	res := []byte(strings.Join(fileContents, "\n"))

	status = "writing lock file"
	j.SendStatus(status)
	err = j.fileWriter.Write(lockFile, res)
	if err != nil {
		cmdErr = util.NewPMJobError(err.Error())
		cmdErr.SetStatus(status)

		return cmdErr
	}

	return nil
}

func (j *Job) runCreateVenvCmd() ([]byte, job.IError) {
	venvName := fmt.Sprintf("%s.venv", filepath.Base(j.GetFile()))
	fpath := filepath.Join(filepath.Dir(j.GetFile()), venvName)
	j.venvPath = fpath

	createVenvCmd, err := j.cmdFactory.MakeCreateVenvCmd(j.venvPath)
	if err != nil {
		cmdErr := util.NewPMJobError(err.Error())
		cmdErr.SetCommand(createVenvCmd.String())

		return nil, cmdErr
	}

	createVenvCmdOutput, err := createVenvCmd.Output()
	if err != nil {
		cmdErr := util.NewPMJobError(j.GetExitError(err).Error())
		cmdErr.SetCommand(createVenvCmd.String())

		return nil, cmdErr
	}

	return createVenvCmdOutput, nil
}

func (j *Job) runInstallCmd() ([]byte, job.IError) {
	var command string
	if j.venvPath != "" {
		binDir := "bin"
		if runtime.GOOS == internalOs.Windows {
			binDir = "Scripts"
		}
		command = filepath.Join(j.venvPath, binDir, pip)
	} else {
		command = pip
	}
	j.pipCommand = command
	installCmd, err := j.cmdFactory.MakeInstallCmd(j.pipCommand, j.GetFile())
	if err != nil {
		cmdErr := util.NewPMJobError(err.Error())
		cmdErr.SetCommand(installCmd.String())

		return nil, cmdErr
	}

	installCmdOutput, err := installCmd.Output()
	if err != nil {
		cmdErr := util.NewPMJobError(j.GetExitError(err).Error())
		cmdErr.SetCommand(installCmd.String())

		return nil, cmdErr
	}

	return installCmdOutput, nil
}

func (j *Job) runCatCmd() ([]byte, job.IError) {
	listCmd, err := j.cmdFactory.MakeCatCmd(j.GetFile())
	if err != nil {
		cmdErr := util.NewPMJobError(err.Error())
		cmdErr.SetCommand(listCmd.String())

		return nil, cmdErr
	}

	listCmdOutput, err := listCmd.Output()
	if err != nil {
		cmdErr := util.NewPMJobError(j.GetExitError(err).Error())
		cmdErr.SetCommand(listCmd.String())

		return nil, cmdErr
	}

	return listCmdOutput, nil
}

func (j *Job) runListCmd() ([]byte, job.IError) {
	listCmd, err := j.cmdFactory.MakeListCmd(j.pipCommand)
	if err != nil {
		cmdErr := util.NewPMJobError(err.Error())
		cmdErr.SetCommand(listCmd.String())

		return nil, cmdErr
	}

	listCmdOutput, err := listCmd.Output()
	if err != nil {
		cmdErr := util.NewPMJobError(j.GetExitError(err).Error())
		cmdErr.SetCommand(listCmd.String())

		return nil, cmdErr
	}

	return listCmdOutput, nil
}

func (j *Job) runShowCmd(packages []string) ([]byte, job.IError) {
	listCmd, err := j.cmdFactory.MakeShowCmd(j.pipCommand, packages)
	if err != nil {
		cmdErr := util.NewPMJobError(err.Error())
		cmdErr.SetCommand(listCmd.String())

		return nil, cmdErr
	}

	listCmdOutput, err := listCmd.Output()
	if err != nil {
		cmdErr := util.NewPMJobError(j.GetExitError(err).Error())
		cmdErr.SetCommand(listCmd.String())

		return nil, cmdErr
	}

	return listCmdOutput, nil
}

func closeFile(j *Job, file *os.File) {
	err := j.fileWriter.Close(file)
	if err != nil {
		jobError := util.NewPMJobError(err.Error())
		j.Errors().Critical(jobError)
	}
}

func (j *Job) parsePipList(pipListOutput string) []string {
	lines := strings.Split(pipListOutput, "\n")
	var packages []string
	for _, line := range lines[2:] {
		fields := strings.Split(line, " ")
		if len(fields) > 0 {
			packages = append(packages, fields[0])
		}
	}

	return packages
}
