package testdata

import (
	"errors"

	"github.com/debricked/cli/pkg/callgraph/config"
	"github.com/debricked/cli/pkg/callgraph/job"
	"github.com/debricked/cli/pkg/callgraph/job/testdata"
	"github.com/debricked/cli/pkg/io/finder"
)

type StrategyMock struct {
	config config.IConfig
	files  []string
	finder finder.IFinder
}

func NewStrategyMock(config config.IConfig, files []string, finder finder.IFinder) StrategyMock {
	return StrategyMock{config, files, finder}
}

func (s StrategyMock) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob
	jobs = append(jobs, testdata.NewJobMock("dir", s.files))

	return jobs, nil
}

type StrategyErrorMock struct {
	config config.IConfig
	files  []string
	finder finder.IFinder
}

func NewStrategyErrorMock(config config.IConfig, files []string, finder finder.IFinder) StrategyErrorMock {
	return StrategyErrorMock{config, files, finder}
}

func (s StrategyErrorMock) Invoke() ([]job.IJob, error) {

	return nil, errors.New("mock-error")
}
