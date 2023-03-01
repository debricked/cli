package tui

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/fatih/color"
)

const (
	title = "Errors"
)

type JobsErrorList struct {
	mirror io.Writer
	jobs   []job.IJob
}

func NewJobsErrorList(mirror io.Writer, jobs []job.IJob) JobsErrorList {
	return JobsErrorList{mirror: mirror, jobs: jobs}
}

func (jobsErrList JobsErrorList) Render() error {
	var listBuffer bytes.Buffer

	formattedTitle := fmt.Sprintf("%s\n", color.BlueString(title))
	underlining := fmt.Sprintf(strings.Repeat("-", len(title)+1) + "\n")
	listBuffer.Write([]byte(formattedTitle))
	listBuffer.Write([]byte(underlining))

	for _, j := range jobsErrList.jobs {
		jobsErrList.addJob(&listBuffer, j)
		listBuffer.Write([]byte("\n"))
	}

	_, err := jobsErrList.mirror.Write(listBuffer.Bytes())

	return err
}

func (jobsErrList JobsErrorList) addJob(list *bytes.Buffer, job job.IJob) {
	var jobString string
	list.Write([]byte(fmt.Sprintf("%s\n", color.YellowString(job.GetFile()))))

	for _, warning := range job.Errors().GetWarningErrors() {
		jobString = fmt.Sprintf("* %s:\n\t%s\n", color.YellowString("Warning"), warning)
		list.Write([]byte(jobString))
	}

	for _, critical := range job.Errors().GetCriticalErrors() {
		jobString = fmt.Sprintf("* %s:\n\t%s\n", color.RedString("Critical"), critical)
		list.Write([]byte(jobString))
	}
}
