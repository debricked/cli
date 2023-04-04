package testdata

import (
	"errors"

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

type StrategyErrorMock struct {
	files []string
}

func NewStrategyErrorMock(files []string) StrategyErrorMock {
	return StrategyErrorMock{files}
}

func (s StrategyErrorMock) Invoke() ([]job.IJob, error) {

	return nil, errors.New("mock-error")
}
