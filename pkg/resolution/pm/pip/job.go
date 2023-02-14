package pip

import (
	"os"
	"strings"

	"github.com/debricked/cli/pkg/resolution/pm/util"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

const (
	fileName = ".debricked-python-dependencies.txt"
)

type Job struct {
	file       string
	cmdFactory ICmdFactory
	fileWriter writer.IFileWriter
	err        error
}

func NewJob(
	file string,
	cmdFactory ICmdFactory,
	fileWriter writer.IFileWriter,
) *Job {
	return &Job{
		file:       file,
		cmdFactory: cmdFactory,
		fileWriter: fileWriter,
	}
}

func (j *Job) File() string {
	return j.file
}

func (j *Job) Error() error {
	return j.err
}

func (j *Job) Run() {

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

	installedPackagesMetadata := string(ShowCmdOutput)

	requiredPackages, err := j.parseRequirements()

	nodes, edges, err := j.parseGraph(requiredPackages, installedPackagesMetadata)

	lockFile, err := j.fileWriter.Create(util.MakePathFromManifestFile(j.file, fileName))
	if err != nil {
		j.err = err
		return
	}
	defer closeFile(j, lockFile)

	var fileContents []string

	fileContents = append(fileContents, nodes...)
	fileContents = append(fileContents, "\n*********\n")
	fileContents = append(fileContents, edges...)

	res := []byte(strings.Join(fileContents, "\n"))

	j.err = j.fileWriter.Write(lockFile, res)
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
