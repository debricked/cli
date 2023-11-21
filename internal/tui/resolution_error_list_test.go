package tui

import (
	"bytes"
	"os"
	"testing"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/job/testdata"
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

	warningErr := job.NewBaseJobError("warning-message")
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
		"|",
		"warning-message\n",
	}
	assertOutput(t, output, contains)
}

func TestRenderCriticalJob(t *testing.T) {
	var listBuffer bytes.Buffer

	warningErr := job.NewBaseJobError("critical-message")
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
		"|",
		"critical-message\n",
	}
	assertOutput(t, output, contains)
}

func TestRenderCriticalAndWarningJob(t *testing.T) {
	var listBuffer bytes.Buffer

	jobMock := testdata.NewJobMock("manifest-file")

	warningErr := job.NewBaseJobError("warning-message")
	jobMock.Errors().Warning(warningErr)

	criticalErr := job.NewBaseJobError("critical-message")
	jobMock.Errors().Critical(criticalErr)

	errList := NewJobsErrorList(&listBuffer, []job.IJob{jobMock})

	err := errList.Render()

	assert.NoError(t, err)
	output := listBuffer.String()
	contains := []string{
		"manifest-file",
		"\n* ",
		"Critical",
		"critical-message\n",
		"Warning",
		"|",
		"warning-message\n",
	}
	assertOutput(t, output, contains)
}

func TestRenderCriticalAndWorkingJob(t *testing.T) {
	var listBuffer bytes.Buffer

	jobWithErrMock := testdata.NewJobMock("manifest-file")

	criticalErr := job.NewBaseJobError("critical-message")
	jobWithErrMock.Errors().Critical(criticalErr)

	jobWorkingMock := testdata.NewJobMock("working-manifest-file")

	errList := NewJobsErrorList(&listBuffer, []job.IJob{jobWithErrMock, jobWorkingMock})

	err := errList.Render()

	assert.NoError(t, err)
	output := listBuffer.String()
	contains := []string{
		"manifest-file",
		"\n* ",
		"Critical",
		"|",
		"critical-message\n",
	}
	assertOutput(t, output, contains)

	assert.NotContains(t, output, jobWorkingMock)
}

func assertOutput(t *testing.T, output string, contains []string) {
	assert.Contains(t, output, "Errors")
	assert.Contains(t, output, "\n-------\n")

	for _, c := range contains {
		assert.Contains(t, output, c)
	}
}
