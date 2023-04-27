package strategy

import (
	"fmt"

	conf "github.com/debricked/cli/pkg/callgraph/config"
	java "github.com/debricked/cli/pkg/callgraph/language/java11"
)

type IFactory interface {
	Make(config conf.IConfig, paths []string) (IStrategy, error)
}

type Factory struct{}

func NewStrategyFactory() Factory {
	return Factory{}
}

func (sf Factory) Make(config conf.IConfig, paths []string) (IStrategy, error) {
	name := config.Language()
	switch name {
	case java.Name:
		return java.NewStrategy(config, paths), nil
	default:
		return nil, fmt.Errorf("failed to make strategy from %s", name)
	}
}
