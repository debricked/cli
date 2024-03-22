package testdata

import (
	"github.com/debricked/cli/internal/callgraph/cgexec"
	"github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/callgraph/finder/testdata"
	"github.com/debricked/cli/internal/callgraph/strategy"
)

type FactoryMock struct{}

func NewStrategyFactoryMock() FactoryMock {
	return FactoryMock{}
}

func (sf FactoryMock) Make(config config.IConfig, paths []string, exclusions []string, ctx cgexec.IContext) (strategy.IStrategy, error) {
	return NewStrategyMock(config, paths, testdata.FinderMock{}, ctx), nil
}

type FactoryErrorMock struct{}

func NewStrategyFactoryErrorMock() FactoryErrorMock {
	return FactoryErrorMock{}
}

func (sf FactoryErrorMock) Make(config config.IConfig, paths []string, exclusions []string, ctx cgexec.IContext) (strategy.IStrategy, error) {
	return NewStrategyErrorMock(config, paths, testdata.FinderMock{}, ctx), nil
}
