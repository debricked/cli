package pip

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/debricked/cli/pkg/resolution/pm/util"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

const (
	lockFileExtension = ".debricked.lock"
	fileName          = ".debricked.lock"
	pip               = "pip"
)

type Job struct {
	file       string
	install    bool
	venvPath   string
	pipCommand string
	cmdFactory ICmdFactory
	fileWriter writer.IFileWriter
	err        error
	status     chan string
}

func NewJob(
	file string,
	install bool,
	cmdFactory ICmdFactory,
	fileWriter writer.IFileWriter,
) *Job {
	return &Job{
		file:       file,
		install:    install,
		cmdFactory: cmdFactory,
		fileWriter: fileWriter,
		status:     make(chan string),
	}
}

func (j *Job) Install() bool {
	return j.install
}

func (j *Job) File() string {
	return j.file
}

func (j *Job) Error() error {
	return j.err
}

func (j *Job) Status() chan string {
	return j.status
}

func (j *Job) Run() {

	if j.install {
		_, err := j.runCreateVenvCmd()

		if err != nil {
			j.err = err

			return
		}

		j.status <- ("created virtualenv for " + j.file + ".venv")

		_, err = j.runInstallCmd()

		if err != nil {
			j.err = err

			return
		}

		j.status <- ("installed requirements in virtualenv for " + j.file + ".venv")

	}

	j.status <- "running cat command"
	catCmdOutput, err := j.runCatCmd()

	if err != nil {
		return
	}

	j.status <- "running list command"
	listCmdOutput, err := j.runListCmd()

	if err != nil {
		return
	}

	j.status <- "running show command"
	installedPackages := j.parsePipList(string(listCmdOutput))
	ShowCmdOutput, err := j.runShowCmd(installedPackages)

	if err != nil {
		return
	}

	j.status <- ("setting up data...")
	lockFileName := "." + filepath.Base(j.file) + lockFileExtension
	lockFile, err := j.fileWriter.Create(util.MakePathFromManifestFile(j.file, lockFileName))

	if err != nil {
		j.err = err

		return
	}
	defer closeFile(j, lockFile)

	var fileContents []string
	fileContents = append(fileContents, string(catCmdOutput))
	fileContents = append(fileContents, "***")
	fileContents = append(fileContents, string(listCmdOutput))
	fileContents = append(fileContents, "***")
	fileContents = append(fileContents, string(ShowCmdOutput))

	res := []byte(strings.Join(fileContents, "\n"))
	j.status <- "writing data..."

	j.err = j.fileWriter.Write(lockFile, res)
}

func (j *Job) runCreateVenvCmd() ([]byte, error) {
	fpath := filepath.Join(filepath.Dir(j.file), filepath.Base(j.file)+".venv")
	j.venvPath = fpath

	createVenvCmd, err := j.cmdFactory.MakeCreateVenvCmd(j.venvPath)
	if err != nil {
		j.err = err

		return nil, err
	}

	createVenvCmdOutput, err := createVenvCmd.Output()
	if err != nil {
		j.err = err

		return nil, err
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
	installCmd, err := j.cmdFactory.MakeInstallCmd(j.pipCommand, j.file)
	if err != nil {
		j.err = err

		return nil, err
	}

	installCmdOutput, err := installCmd.Output()
	if err != nil {
		j.err = err

		return nil, err
	}

	return installCmdOutput, nil
}

func (j *Job) runCatCmd() ([]byte, error) {
	listCmd, err := j.cmdFactory.MakeCatCmd(j.file)
	if err != nil {
		j.err = err

		return nil, err
	}

	listCmdOutput, err := listCmd.Output()
	if err != nil {
		j.err = err

		return nil, err
	}

	return listCmdOutput, nil
}

func (j *Job) runListCmd() ([]byte, error) {
	listCmd, err := j.cmdFactory.MakeListCmd(j.pipCommand)
	if err != nil {
		j.err = err

		return nil, err
	}

	listCmdOutput, err := listCmd.Output()
	if err != nil {
		j.err = err

		return nil, err
	}

	return listCmdOutput, nil
}

func (j *Job) runShowCmd(packages []string) ([]byte, error) {
	listCmd, err := j.cmdFactory.MakeShowCmd(j.pipCommand, packages)
	if err != nil {
		j.err = err

		return nil, err
	}

	listCmdOutput, err := listCmd.Output()
	if err != nil {
		j.err = err

		return nil, err
	}

	return listCmdOutput, nil
}

func closeFile(job *Job, file *os.File) {
	err := job.fileWriter.Close(file)
	if err != nil {
		job.err = err
	}
}

func (j *Job) parsePipList(pipListOutput string) []string {
	lines := strings.Split(pipListOutput, "\n")
	packages := []string{}
	for _, line := range lines[2:] {
		fields := strings.Split(line, " ")
		if len(fields) > 0 {
			packages = append(packages, fields[0])
		}
	}

	return packages
}
