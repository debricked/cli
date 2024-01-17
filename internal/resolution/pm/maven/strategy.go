package maven

import (
	"github.com/debricked/cli/internal/resolution/job"
)

type Strategy struct {
	files      []string
	cmdFactory ICmdFactory
}

func NewStrategy(files []string) Strategy {
	return Strategy{files, CmdFactory{}}
}

func (s Strategy) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob

	for _, file := range s.files {
		jobs = append(jobs, NewJob(file, s.cmdFactory, PomService{}))
	}

	return jobs, nil
}
