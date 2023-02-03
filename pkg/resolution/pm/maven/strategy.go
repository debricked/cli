package maven

import (
	"github.com/debricked/cli/pkg/resolution/job"
)

type Strategy struct {
	files      []string
	cmdFactory ICmdFactory
}

func NewStrategy(files []string) Strategy {
	return Strategy{files, CmdFactory{}}
}

func (s Strategy) Invoke() []job.IJob {
	var jobs []job.IJob
	for _, file := range s.files {
		jobs = append(jobs, NewJob(file, s.cmdFactory))
	}

	return jobs
}
