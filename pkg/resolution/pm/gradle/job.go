package gradle

import (
	"os"

	"github.com/debricked/cli/pkg/resolution/pm/util"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

const (
	fileName = ".debricked-gradle-dependencies.txt"
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
	dependenciesCmd, err := j.cmdFactory.MakeDependenciesCmd()
	if err != nil {
		j.err = err

		return
	}

	output, err := dependenciesCmd.Output()
	if err != nil {
		j.err = err

		return
	}

	lockFile, err := j.fileWriter.Create(util.MakePathFromManifestFile(j.file, fileName))
	if err != nil {
		j.err = err

		return
	}
	defer closeFile(j, lockFile)

	j.err = j.fileWriter.Write(lockFile, output)
}

func closeFile(job *Job, file *os.File) {
	err := job.fileWriter.Close(file)
	if err != nil {
		job.err = err
	}
}
