package resolution

import (
	"testing"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewResolution(t *testing.T) {
	res := NewResolution(nil)
	assert.NotNil(t, res)

	res = NewResolution([]job.IJob{})
	assert.NotNil(t, res)

	res = NewResolution([]job.IJob{testdata.NewJobMock("")})
	assert.NotNil(t, res)

	res = NewResolution([]job.IJob{testdata.NewJobMock(""), testdata.NewJobMock("")})
	assert.NotNil(t, res)
}

func TestJobs(t *testing.T) {
	res := NewResolution(nil)
	assert.Empty(t, res.Jobs())

	res.jobs = []job.IJob{}
	assert.Len(t, res.Jobs(), 0)

	res.jobs = []job.IJob{testdata.NewJobMock("")}
	assert.Len(t, res.Jobs(), 1)

	res.jobs = []job.IJob{testdata.NewJobMock(""), testdata.NewJobMock("")}
	assert.Len(t, res.Jobs(), 2)
}

func TestHasError(t *testing.T) {
	res := NewResolution(nil)
	assert.False(t, res.HasErr())

	res.jobs = []job.IJob{testdata.NewJobMock("")}
	assert.False(t, res.HasErr())

	jobMock := testdata.NewJobMock("")
	jobMock.SetErr(job.NewBaseJobError("error"))
	res.jobs = append(res.jobs, jobMock)
	assert.True(t, res.HasErr())
}
