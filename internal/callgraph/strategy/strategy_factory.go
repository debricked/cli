package strategy

import (
	"fmt"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	conf "github.com/debricked/cli/internal/callgraph/config"
	golanfinder "github.com/debricked/cli/internal/callgraph/finder/golangfinder"
	"github.com/debricked/cli/internal/callgraph/finder/javafinder"
	"github.com/debricked/cli/internal/callgraph/language/golang"
	java "github.com/debricked/cli/internal/callgraph/language/java11"
)

type IFactory interface {
	Make(config conf.IConfig, paths []string, exclusions []string, ctx cgexec.IContext) (IStrategy, error)
}

type Factory struct{}

func NewStrategyFactory() Factory {
	return Factory{}
}

func (sf Factory) Make(config conf.IConfig, paths []string, exclusions []string, ctx cgexec.IContext) (IStrategy, error) {
	name := config.Language()
	switch name {
	case java.Name:
		return java.NewStrategy(config, paths, exclusions, javafinder.JavaFinder{}, ctx), nil
	case golang.Name:
		return golang.NewStrategy(config, paths, exclusions, golanfinder.GolangFinder{}, ctx), nil
	default:
		return nil, fmt.Errorf("failed to make strategy from %s", name)
	}
}
