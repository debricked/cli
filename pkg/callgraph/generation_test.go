package callgraph

import (
	"errors"
	"testing"

	"github.com/debricked/cli/pkg/callgraph/job"
	"github.com/debricked/cli/pkg/callgraph/job/testdata"
	"github.com/stretchr/testify/assert"
)

const testDir = "dir"

var testFiles = []string{"file"}

func TestNewGeneration(t *testing.T) {
	res := NewGeneration(nil)
	assert.NotNil(t, res)

	res = NewGeneration([]job.IJob{})
	assert.NotNil(t, res)

	res = NewGeneration([]job.IJob{testdata.NewJobMock(testDir, testFiles)})
	assert.NotNil(t, res)

	res = NewGeneration([]job.IJob{testdata.NewJobMock(testDir, testFiles), testdata.NewJobMock(testDir, testFiles)})
	assert.NotNil(t, res)
}

func TestJobs(t *testing.T) {
	res := NewGeneration(nil)
	assert.Empty(t, res.Jobs())

	res.jobs = []job.IJob{}
	assert.Len(t, res.Jobs(), 0)

	res.jobs = []job.IJob{testdata.NewJobMock(testDir, testFiles)}
	assert.Len(t, res.Jobs(), 1)

	res.jobs = []job.IJob{testdata.NewJobMock(testDir, testFiles), testdata.NewJobMock(testDir, testFiles)}
	assert.Len(t, res.Jobs(), 2)
}

func TestHasError(t *testing.T) {
	res := NewGeneration(nil)
	assert.False(t, res.HasErr())

	res.jobs = []job.IJob{testdata.NewJobMock(testDir, testFiles)}
	assert.False(t, res.HasErr())

	jobMock := testdata.NewJobMock(testDir, testFiles)
	jobMock.SetErr(errors.New("error"))
	res.jobs = append(res.jobs, jobMock)
	assert.True(t, res.HasErr())
}
