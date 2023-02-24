package resolution

import (
	"fmt"

	"github.com/debricked/cli/pkg/resolution/job"
)

type IScheduler interface {
	Schedule(jobs []job.IJob) (IResolution, error)
}

type Scheduler struct{}

func NewScheduler() Scheduler {
	return Scheduler{}
}

func (s Scheduler) Schedule(jobs []job.IJob) (IResolution, error) {
	for _, j := range jobs {
		j.Run()
		if j.Error() != nil {
			fmt.Println(j.Error())
		}
	}

	return NewResolution(jobs), nil
}
