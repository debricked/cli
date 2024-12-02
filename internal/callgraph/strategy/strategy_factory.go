package strategy

import (
	"fmt"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	conf "github.com/debricked/cli/internal/callgraph/config"
	golangfinder "github.com/debricked/cli/internal/callgraph/finder/golangfinder"
	"github.com/debricked/cli/internal/callgraph/finder/javafinder"
	"github.com/debricked/cli/internal/callgraph/language/golang"
	"github.com/debricked/cli/internal/callgraph/language/java"
)

type IFactory interface {
	Make(config conf.IConfig, paths []string, exclusions []string, inclusions []string, ctx cgexec.IContext) (IStrategy, error)
}

type Factory struct{}

func NewStrategyFactory() Factory {
	return Factory{}
}

func (sf Factory) Make(
	config conf.IConfig,
	paths []string,
	exclusions []string,
	inclusions []string,
	ctx cgexec.IContext,
) (IStrategy, error) {
	name := config.Language()
	switch name {
	case java.Name:
		return java.NewStrategy(config, paths, exclusions, inclusions, javafinder.JavaFinder{}, ctx), nil
	case golang.Name:
		return golang.NewStrategy(config, paths, exclusions, inclusions, golangfinder.GolangFinder{}, ctx), nil
	default:
		return nil, fmt.Errorf("failed to make strategy from %s", name)
	}
}
