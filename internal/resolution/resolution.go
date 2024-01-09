package resolution

import "github.com/debricked/cli/internal/resolution/job"

type IResolution interface {
	Jobs() []job.IJob
	HasErr() bool
	GetJobErrorCount() int
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

func (r Resolution) GetJobErrorCount() int {
	count := 0
	for _, j := range r.Jobs() {
		if j.Errors().HasError() {
			count++
		}
	}

	return count
}
