package tui

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/job/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJobsErrorList(t *testing.T) {
	mirror := os.Stdout
	errList := NewJobsErrorList(mirror, []job.IJob{})
	assert.NotNil(t, errList)
}

func TestRenderNoJobs(t *testing.T) {
	var listBuffer bytes.Buffer
	errList := NewJobsErrorList(&listBuffer, []job.IJob{})

	err := errList.Render()

	assert.NoError(t, err)
	output := listBuffer.String()
	assertOutput(t, output, nil)
}

func TestRenderWarningJob(t *testing.T) {
	var listBuffer bytes.Buffer

	warningErr := errors.New("warning-message")
	jobMock := testdata.NewJobMock("file")
	jobMock.Errors().Warning(warningErr)
	errList := NewJobsErrorList(&listBuffer, []job.IJob{jobMock})

	err := errList.Render()

	assert.NoError(t, err)
	output := listBuffer.String()
	contains := []string{
		"file",
		"\n* ",
		"Warning",
		":\n\twarning-message\n",
	}
	assertOutput(t, output, contains)
}

func TestRenderCriticalJob(t *testing.T) {
	var listBuffer bytes.Buffer

	warningErr := errors.New("critical-message")
	jobMock := testdata.NewJobMock("file")
	jobMock.Errors().Critical(warningErr)
	errList := NewJobsErrorList(&listBuffer, []job.IJob{jobMock})

	err := errList.Render()

	assert.NoError(t, err)
	output := listBuffer.String()
	contains := []string{
		"file",
		"\n* ",
		"Critical",
		":\n\tcritical-message\n",
	}
	assertOutput(t, output, contains)
}

func TestRenderCriticalAndWarningJob(t *testing.T) {
	var listBuffer bytes.Buffer

	jobMock := testdata.NewJobMock("file")

	warningErr := errors.New("warning-message")
	jobMock.Errors().Warning(warningErr)

	criticalErr := errors.New("critical-message")
	jobMock.Errors().Critical(criticalErr)

	errList := NewJobsErrorList(&listBuffer, []job.IJob{jobMock})

	err := errList.Render()

	assert.NoError(t, err)
	output := listBuffer.String()
	contains := []string{
		"file",
		"\n* ",
		"Critical",
		":\n\tcritical-message\n",
		"Warning",
		":\n\twarning-message\n",
	}
	assertOutput(t, output, contains)
	fmt.Println(output)
}

func assertOutput(t *testing.T, output string, contains []string) {
	assert.Contains(t, output, "Errors")
	assert.Contains(t, output, "\n-------\n")

	for _, c := range contains {
		assert.Contains(t, output, c)
	}
}
