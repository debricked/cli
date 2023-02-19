package pip

import (
	"os"
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
		_, err := j.runInstallCmd()
		if err != nil {
			j.err = err
			return
		}
		// TODO if unable to install (many possible issues)
		// then we should parse and let the user know what went wrong on installation

		// TODO create separate env for installation as to not collide with other installed packages
		// TODO install in a virtualenv
	}

	catCmdOutput, err := j.runCatCmd()
	if err != nil {
		return
	}

	listCmdOutput, err := j.runListCmd()
	if err != nil {
		return
	}

	installedPackages, err := j.parsePipList(string(listCmdOutput))
	if err != nil {
		return
	}

	ShowCmdOutput, err := j.runShowCmd(installedPackages)
	if err != nil {
		return
	}

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

	j.err = j.fileWriter.Write(lockFile, res)
}

func (j *Job) runCatCmd() ([]byte, error) {
	catCmd, err := j.cmdFactory.MakeCatCmd(j.file)
	if err != nil {
		j.err = err

		return nil, err
	}

	catCmdOutput, err := catCmd.Output()
	if err != nil {
		j.err = err

		return nil, err
	}

	return catCmdOutput, nil
}

func (j *Job) runListCmd() ([]byte, error) {
	listCmd, err := j.cmdFactory.MakeListCmd()
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

func (j *Job) runInstallCmd() ([]byte, error) {
	installCmd, err := j.cmdFactory.MakeInstallCmd(j.file)
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

func (j *Job) runShowCmd(packages []string) ([]byte, error) {
	listCmd, err := j.cmdFactory.MakeShowCmd(packages)
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
