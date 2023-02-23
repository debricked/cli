package pip

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/debricked/cli/pkg/resolution/pm/util"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

const (
	fileName = ".debricked-pip-tree.txt"
)

type Job struct {
	file       string
	install    bool
	venvPath   string
	pipCommand string
	cmdFactory ICmdFactory
	fileWriter writer.IFileWriter
	err        error
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

func (j *Job) Run() {

	if j.install {
		_, err := j.runCreateVenvCmd()

		if err != nil {
			j.err = err
			return
		}

		fmt.Println("Created virtualenv for " + j.file + ".venv")

		_, err = j.runActivateVenvCmd()

		if err != nil {
			j.err = err
			fmt.Println(err)
			return
		}

		fmt.Println("Activated virtualenv for " + j.file + ".venv")

		_, err = j.runInstallCmd()

		if err != nil {
			j.err = err
			return
		}

		fmt.Println("Installed requirements in virtualenv for " + j.file + ".venv")
		// TODO if unable to install (many possible issues)
		// then we should parse and let the user know what went wrong on installation

	}

	fmt.Println("Running cat command")
	catCmdOutput, err := j.runCatCmd()

	if err != nil {
		return
	}

	fmt.Println("Running list command")
	listCmdOutput, err := j.runListCmd()

	if err != nil {
		return
	}

	fmt.Println("Running show command")
	installedPackages := j.parsePipList(string(listCmdOutput))
	ShowCmdOutput, err := j.runShowCmd(installedPackages)

	if err != nil {
		return
	}

	fmt.Println("Setting up data...")
	lockFile, err := j.fileWriter.Create(util.MakePathFromManifestFile(j.file, fileName))

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
	fmt.Println("Writing data...")

	j.err = j.fileWriter.Write(lockFile, res)
}

func (j *Job) runCreateVenvCmd() ([]byte, error) {
	fpath := filepath.Join(filepath.Dir(j.file), filepath.Base(j.file)+".venv")
	j.venvPath = fpath
	fmt.Println("MakeCreateVenvCmd", j.venvPath)

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

func (j *Job) runActivateVenvCmd() ([]byte, error) {
	activateVenvCmd, err := j.cmdFactory.MakeActivateVenvCmd(j.file)
	if err != nil {
		j.err = err
		return nil, err
	}

	activateVenvCmdOutput, err := activateVenvCmd.Output()
	if err != nil {
		j.err = err
		return nil, err
	}

	return activateVenvCmdOutput, nil
}

func (j *Job) runInstallCmd() ([]byte, error) {
	var command string
	if j.venvPath != "" {
		command = filepath.Join(j.venvPath, "bin", "pip")
	} else {
		command = "pip"
	}
	fmt.Println("MakeInstallCmd", command, j.venvPath)
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
