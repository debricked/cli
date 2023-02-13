package pip

import (
	"bufio"
	"fmt"
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

	file, err := os.Open(j.file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	packages := []string{}

	// Loop through each line of the file
	for scanner.Scan() {
		line := scanner.Text()

		// Split the line by spaces
		fields := strings.Fields(line)

		// If the line starts with a # symbol, it's a comment, so we skip it
		if len(fields) == 0 || fields[0][0] == '#' {
			continue
		}

		// Print the first field, which should be the package name
		packages = append(packages, fields[0])
	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(packages)

	fmt.Println("Run List cmd")
	listCmdOutput, err := j.runListCmd()
	if err != nil {
		return
	}
	fmt.Println("Make file")
	lockFile, err := j.fileWriter.Create(util.MakePathFromManifestFile(j.file, fileName))
	if err != nil {
		j.err = err
		return
	}
	defer closeFile(j, lockFile)

	var fileContents []byte
	fileContents = append(fileContents, []byte("\nVERYGOOD DELIMITER\n")...)
	fileContents = append(fileContents, listCmdOutput...)

	j.err = j.fileWriter.Write(lockFile, fileContents)
}

func (j *Job) parseRequiredNamesCmd() ([]byte, error) {
	return nil, nil
}

func (j *Job) parseInstalledNamesCmd() ([]byte, error) {
	return nil, nil
}

func (j *Job) parseRelationNamesCmd() ([]byte, error) {
	return nil, nil
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

func closeFile(job *Job, file *os.File) {
	err := job.fileWriter.Close(file)
	if err != nil {
		job.err = err
	}
}
