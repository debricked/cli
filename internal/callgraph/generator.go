package callgraph

import (
	"os"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	"github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/callgraph/finder"
	"github.com/debricked/cli/internal/callgraph/job"
	"github.com/debricked/cli/internal/callgraph/strategy"
	"github.com/debricked/cli/internal/tui"
)

type IGenerator interface {
	GenerateWithTimer(paths []string, exclusions []string, configs []config.IConfig, timeout int) error
	Generate(paths []string, exclusions []string, configs []config.IConfig, ctx cgexec.IContext) error
}

type Generator struct {
	finder          finder.IFinder
	strategyFactory strategy.IFactory
	scheduler       IScheduler
	Generation      IGeneration
}

func NewGenerator(
	finder finder.IFinder,
	strategyFactory strategy.IFactory,
	scheduler IScheduler,
) *Generator {
	return &Generator{
		finder,
		strategyFactory,
		scheduler,
		Generation{},
	}
}

func (g *Generator) GenerateWithTimer(paths []string, exclusions []string, configs []config.IConfig, timeout int) error {
	result := make(chan error)
	ctx, cancel := cgexec.NewContext(timeout)
	defer cancel()

	go func() {
		result <- g.Generate(paths, exclusions, configs, &ctx)
	}()

	// Wait for the result or timeout
	err := <-result

	return err
}

func (g *Generator) Generate(paths []string, exclusions []string, configs []config.IConfig, ctx cgexec.IContext) error {
	targetPath := ".debrickedTmpFolder"
	debrickedExclusions := []string{targetPath}
	exclusions = append(exclusions, debrickedExclusions...)
	files, _ := g.finder.FindFiles(paths, exclusions)

	var jobs []job.IJob
	for _, config := range configs {
		s, strategyErr := g.strategyFactory.Make(config, files, paths, exclusions, g.finder, ctx)
		if strategyErr == nil {
			newJobs, err := s.Invoke()
			if err != nil {
				return err
			}
			jobs = append(jobs, newJobs...)
		}
	}

	generation, err := g.scheduler.Schedule(jobs, ctx)
	g.Generation = generation

	if generation.HasErr() {
		jobErrList := tui.NewCallgraphJobsErrorList(os.Stdout, generation.Jobs())
		err = jobErrList.Render()
	}

	return err
}
