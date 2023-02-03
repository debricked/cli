package resolution

import (
	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/job/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
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
	s := NewScheduler()
	assert.NotNil(t, s)
}

func TestSchedule(t *testing.T) {
	s := NewScheduler()
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
}
