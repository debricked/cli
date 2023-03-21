package testdata

import (
	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/job/testdata"
)

type StrategyMock struct {
	files []string
}

func NewStrategyMock(files []string) StrategyMock {
	return StrategyMock{files}
}

func (s StrategyMock) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob
	for _, file := range s.files {
		jobs = append(jobs, testdata.NewJobMock(file))
	}

	return jobs, nil
}
