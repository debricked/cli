package resolution

import (
	"errors"
	"testing"

	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/job/testdata"
	"github.com/stretchr/testify/assert"
)

type SchedulerMock struct {
	Err error
}

func (s SchedulerMock) Schedule(jobs []job.IJob) (IResolution, error) {
	for _, j := range jobs {
		j.Run()
	}

	return NewResolution(jobs), s.Err
}

func TestNewScheduler(t *testing.T) {
	s := NewScheduler(10)
	assert.NotNil(t, s)
}

func TestSchedule(t *testing.T) {
	s := NewScheduler(10)
	res, err := s.Schedule([]job.IJob{testdata.NewJobMock("")})
	assert.NoError(t, err)
	assert.Len(t, res.Jobs(), 1)

	res, err = s.Schedule([]job.IJob{})
	assert.NoError(t, err)
	assert.Len(t, res.Jobs(), 0)

	res, err = s.Schedule(nil)
	assert.NoError(t, err)
	assert.Len(t, res.Jobs(), 0)

	res, err = s.Schedule([]job.IJob{testdata.NewJobMock(""), testdata.NewJobMock("")})
	assert.NoError(t, err)
	assert.Len(t, res.Jobs(), 2)
	for _, j := range res.Jobs() {
		assert.False(t, j.Errors().HasError())
	}
}

func TestScheduleJobErr(t *testing.T) {
	s := NewScheduler(10)
	jobMock := testdata.NewJobMock("")
	jobErr := errors.New("job-error")
	jobMock.SetErr(jobErr)
	res, err := s.Schedule([]job.IJob{jobMock})
	assert.NoError(t, err)
	assert.Len(t, res.Jobs(), 1)
	j := res.Jobs()[0]
	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), jobErr)
}
