package tui

import (
	"bytes"
	"errors"
	"os"
	"testing"

	"github.com/debricked/cli/internal/callgraph/job"
	"github.com/debricked/cli/internal/callgraph/job/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewCallgraphJobsErrorList(t *testing.T) {
	mirror := os.Stdout
	errList := NewCallgraphJobsErrorList(mirror, []job.IJob{})
	assert.NotNil(t, errList)
}

func TestRenderNoCallgraphJobs(t *testing.T) {
	var listBuffer bytes.Buffer
	errList := NewCallgraphJobsErrorList(&listBuffer, []job.IJob{})

	err := errList.Render()

	assert.NoError(t, err)
	output := listBuffer.String()
	assertOutput(t, output, nil)
}

func TestRenderWarningCallgraphJob(t *testing.T) {
	var listBuffer bytes.Buffer

	warningErr := errors.New("warning-message")
	jobMock := testdata.NewJobMock("file", nil)
	jobMock.Errors().Warning(warningErr)
	errList := NewCallgraphJobsErrorList(&listBuffer, []job.IJob{jobMock})

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

func TestRenderCriticalCallgraphJob(t *testing.T) {
	var listBuffer bytes.Buffer

	warningErr := errors.New("critical-message")
	jobMock := testdata.NewJobMock("file", nil)
	jobMock.Errors().Critical(warningErr)
	errList := NewCallgraphJobsErrorList(&listBuffer, []job.IJob{jobMock})

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

func TestRenderCriticalAndWarningCallgraphJob(t *testing.T) {
	var listBuffer bytes.Buffer

	jobMock := testdata.NewJobMock("manifest-file", nil)

	warningErr := errors.New("warning-message")
	jobMock.Errors().Warning(warningErr)

	criticalErr := errors.New("critical-message")
	jobMock.Errors().Critical(criticalErr)

	errList := NewCallgraphJobsErrorList(&listBuffer, []job.IJob{jobMock})

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

func TestRenderCriticalAndWorkingCallgraphJob(t *testing.T) {
	var listBuffer bytes.Buffer

	jobWithErrMock := testdata.NewJobMock("manifest-file", nil)

	criticalErr := errors.New("critical-message")
	jobWithErrMock.Errors().Critical(criticalErr)

	jobWorkingMock := testdata.NewJobMock("working-manifest-file", nil)

	errList := NewCallgraphJobsErrorList(&listBuffer, []job.IJob{jobWithErrMock, jobWorkingMock})

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
