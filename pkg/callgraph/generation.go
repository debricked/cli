package callgraph

import "github.com/debricked/cli/pkg/callgraph/job"

type IGeneration interface {
	Jobs() []job.IJob
	HasErr() bool
}

type Generation struct {
	jobs []job.IJob
}

func NewGeneration(jobs []job.IJob) Generation {
	return Generation{jobs}
}

func (g Generation) Jobs() []job.IJob {
	return g.jobs
}

func (g Generation) HasErr() bool {
	for _, j := range g.Jobs() {
		if j.Errors().HasError() {
			return true
		}
	}

	return false
}
