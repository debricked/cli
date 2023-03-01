package resolution

import "github.com/debricked/cli/pkg/resolution/job"

type IResolution interface {
	Jobs() []job.IJob
	HasErr() bool
}

type Resolution struct {
	jobs []job.IJob
}

func NewResolution(jobs []job.IJob) Resolution {
	return Resolution{jobs}
}

func (r Resolution) Jobs() []job.IJob {
	return r.jobs
}

func (r Resolution) HasErr() bool {
	for _, j := range r.Jobs() {
		if j.Errors().HasError() {
			return true
		}
	}

	return false
}
