package strategy

import (
	"fmt"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	conf "github.com/debricked/cli/internal/callgraph/config"
	java "github.com/debricked/cli/internal/callgraph/language/java11"
	"github.com/debricked/cli/internal/io/finder"
)

type IFactory interface {
	Make(config conf.IConfig, paths []string, finder finder.IFinder, ctx cgexec.IContext) (IStrategy, error)
}

type Factory struct{}

func NewStrategyFactory() Factory {
	return Factory{}
}

func (sf Factory) Make(config conf.IConfig, paths []string, finder finder.IFinder, ctx cgexec.IContext) (IStrategy, error) {
	name := config.Language()
	switch name {
	case java.Name:
		return java.NewStrategy(config, paths, finder, ctx), nil
	default:
		return nil, fmt.Errorf("failed to make strategy from %s", name)
	}
}
