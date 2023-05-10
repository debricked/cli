package strategy

import (
	"fmt"

	conf "github.com/debricked/cli/pkg/callgraph/config"
	java "github.com/debricked/cli/pkg/callgraph/language/java11"
	"github.com/debricked/cli/pkg/io/finder"
)

type IFactory interface {
	Make(config conf.IConfig, paths []string, finder finder.IFinder) (IStrategy, error)
}

type Factory struct{}

func NewStrategyFactory() Factory {
	return Factory{}
}

func (sf Factory) Make(config conf.IConfig, paths []string, finder finder.IFinder) (IStrategy, error) {
	name := config.Language()
	switch name {
	case java.Name:
		return java.NewStrategy(config, paths, finder), nil
	default:
		return nil, fmt.Errorf("failed to make strategy from %s", name)
	}
}
