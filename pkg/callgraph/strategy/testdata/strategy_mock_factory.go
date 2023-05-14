package testdata

import (
	"github.com/debricked/cli/pkg/callgraph/config"
	"github.com/debricked/cli/pkg/callgraph/strategy"
	"github.com/debricked/cli/pkg/io/finder"
)

type FactoryMock struct{}

func NewStrategyFactoryMock() FactoryMock {
	return FactoryMock{}
}

func (sf FactoryMock) Make(config config.IConfig, paths []string, finder finder.IFinder) (strategy.IStrategy, error) {

	return NewStrategyMock(config, paths, finder), nil
}

type FactoryErrorMock struct{}

func NewStrategyFactoryErrorMock() FactoryErrorMock {
	return FactoryErrorMock{}
}

func (sf FactoryErrorMock) Make(config config.IConfig, paths []string, finder finder.IFinder) (strategy.IStrategy, error) {
	return NewStrategyErrorMock(config, paths, finder), nil
}
