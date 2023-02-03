package resolution

import "github.com/debricked/cli/pkg/resolution/job"

type IResolution interface {
	Jobs() []job.IJob
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
