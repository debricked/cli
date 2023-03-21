package gomod

import (
	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

type Strategy struct {
	files []string
}

func (s Strategy) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob
	for _, file := range s.files {
		jobs = append(jobs, NewJob(file, CmdFactory{}, writer.FileWriter{}))
	}

	return jobs, nil
}

func NewStrategy(files []string) Strategy {
	return Strategy{files}
}
