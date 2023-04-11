package testdata

import (
	"github.com/debricked/cli/pkg/resolution/file"
	"github.com/debricked/cli/pkg/resolution/strategy"
)

type FactoryMock struct{}

func NewStrategyFactoryMock() FactoryMock {
	return FactoryMock{}
}

func (sf FactoryMock) Make(pmFileBatch file.IBatch, paths []string) (strategy.IStrategy, error) {

	return NewStrategyMock(pmFileBatch.Files()), nil
}

type FactoryErrorMock struct{}

func NewStrategyFactoryErrorMock() FactoryErrorMock {
	return FactoryErrorMock{}
}

func (sf FactoryErrorMock) Make(pmFileBatch file.IBatch, paths []string) (strategy.IStrategy, error) {

	return NewStrategyErrorMock(pmFileBatch.Files()), nil
}
