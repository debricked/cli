package testdata

import (
	"errors"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	"github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/callgraph/finder"
	"github.com/debricked/cli/internal/callgraph/job"
	"github.com/debricked/cli/internal/callgraph/job/testdata"
)

type StrategyMock struct {
	config config.IConfig
	files  []string
	finder finder.IFinder
	ctx    cgexec.IContext
}

func NewStrategyMock(config config.IConfig, files []string, finder finder.IFinder, ctx cgexec.IContext) StrategyMock {
	return StrategyMock{config, files, finder, ctx}
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
	ctx    cgexec.IContext
}

func NewStrategyErrorMock(config config.IConfig, files []string, finder finder.IFinder, ctx cgexec.IContext) StrategyErrorMock {
	return StrategyErrorMock{config, files, finder, ctx}
}

func (s StrategyErrorMock) Invoke() ([]job.IJob, error) {

	return nil, errors.New("mock-error")
}
