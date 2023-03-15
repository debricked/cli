package pip

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/pm/util"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

const (
	lockFileExtension = ".debricked.lock"
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
}

func NewJob(
	file string,
	install bool,
	cmdFactory ICmdFactory,
	fileWriter writer.IFileWriter,
) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(file),
		install:    install,
		cmdFactory: cmdFactory,
		fileWriter: fileWriter,
	}
}

func (j *Job) Install() bool {
	return j.install
}

func (j *Job) Run() {
	if j.install {
		j.SendStatus(fmt.Sprintf("creating virtualenv for %s.venv", filepath.Base(j.GetFile())))
		_, err := j.runCreateVenvCmd()
		if err != nil {
			j.Errors().Critical(err)

			return
		}

		j.SendStatus(fmt.Sprintf("installing requirements in virtualenv for %s.venv", filepath.Base(j.GetFile())))
		_, err = j.runInstallCmd()
		if err != nil {
			j.Errors().Critical(err)

			return
		}
	}

	j.SendStatus("running cat command")
	catCmdOutput, err := j.runCatCmd()
	if err != nil {
		j.Errors().Critical(err)

		return
	}

	j.SendStatus("running list command")
	listCmdOutput, err := j.runListCmd()
	if err != nil {
		j.Errors().Critical(err)

		return
	}

	j.SendStatus("running show command")
	installedPackages := j.parsePipList(string(listCmdOutput))
	ShowCmdOutput, err := j.runShowCmd(installedPackages)
	if err != nil {
		j.Errors().Critical(err)

		return
	}

	j.SendStatus("setting up data...")
	lockFileName := fmt.Sprintf(".%s%s", filepath.Base(j.GetFile()), lockFileExtension)
	lockFile, err := j.fileWriter.Create(util.MakePathFromManifestFile(j.GetFile(), lockFileName))
	if err != nil {
		j.Errors().Critical(err)

		return
	}
	defer closeFile(j, lockFile)

	var fileContents []string
	fileContents = append(fileContents, string(catCmdOutput))
	fileContents = append(fileContents, lockFileDelimiter)
	fileContents = append(fileContents, string(listCmdOutput))
	fileContents = append(fileContents, lockFileDelimiter)
	fileContents = append(fileContents, string(ShowCmdOutput))

	res := []byte(strings.Join(fileContents, "\n"))
	j.SendStatus("writing data...")

	err = j.fileWriter.Write(lockFile, res)
	if err != nil {
		j.Errors().Critical(err)
	}
}

func (j *Job) runCreateVenvCmd() ([]byte, error) {
	venvName := fmt.Sprintf("%s.venv", filepath.Base(j.GetFile()))
	fpath := filepath.Join(filepath.Dir(j.GetFile()), venvName)
	j.venvPath = fpath

	createVenvCmd, err := j.cmdFactory.MakeCreateVenvCmd(j.venvPath)
	if err != nil {
		return nil, err
	}

	createVenvCmdOutput, err := createVenvCmd.Output()
	if err != nil {
		return nil, j.GetExitError(err)
	}

	return createVenvCmdOutput, nil
}

func (j *Job) runInstallCmd() ([]byte, error) {
	var command string
	if j.venvPath != "" {
		command = filepath.Join(j.venvPath, "bin", pip)
	} else {
		command = pip
	}
	j.pipCommand = command
	installCmd, err := j.cmdFactory.MakeInstallCmd(j.pipCommand, j.GetFile())
	if err != nil {
		return nil, err
	}

	installCmdOutput, err := installCmd.Output()
	if err != nil {
		return nil, j.GetExitError(err)
	}

	return installCmdOutput, nil
}

func (j *Job) runCatCmd() ([]byte, error) {
	listCmd, err := j.cmdFactory.MakeCatCmd(j.GetFile())
	if err != nil {
		return nil, err
	}

	listCmdOutput, err := listCmd.Output()
	if err != nil {
		return nil, j.GetExitError(err)
	}

	return listCmdOutput, nil
}

func (j *Job) runListCmd() ([]byte, error) {
	listCmd, err := j.cmdFactory.MakeListCmd(j.pipCommand)
	if err != nil {
		return nil, err
	}

	listCmdOutput, err := listCmd.Output()
	if err != nil {
		return nil, j.GetExitError(err)
	}

	return listCmdOutput, nil
}

func (j *Job) runShowCmd(packages []string) ([]byte, error) {
	listCmd, err := j.cmdFactory.MakeShowCmd(j.pipCommand, packages)
	if err != nil {
		return nil, err
	}

	listCmdOutput, err := listCmd.Output()
	if err != nil {
		return nil, j.GetExitError(err)
	}

	return listCmdOutput, nil
}

func closeFile(job *Job, file *os.File) {
	err := job.fileWriter.Close(file)
	if err != nil {
		job.Errors().Critical(err)
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
