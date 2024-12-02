package callgraph

import (
	"os"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	"github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/callgraph/job"
	"github.com/debricked/cli/internal/callgraph/strategy"
	"github.com/debricked/cli/internal/tui"
)

type DebrickedOptions struct {
	Paths      []string
	Exclusions []string
	Inclusions []string
	Configs    []config.IConfig
	Timeout    int
	Version    string
}

type IGenerator interface {
	GenerateWithTimer(options DebrickedOptions) error
	Generate(options DebrickedOptions, ctx cgexec.IContext) error
}

type Generator struct {
	strategyFactory strategy.IFactory
	scheduler       IScheduler
	Generation      IGeneration
}

func NewGenerator(
	strategyFactory strategy.IFactory,
	scheduler IScheduler,
) *Generator {
	return &Generator{
		strategyFactory,
		scheduler,
		Generation{},
	}
}

func (g *Generator) GenerateWithTimer(options DebrickedOptions) error {
	result := make(chan error)
	ctx, cancel := cgexec.NewContext(options.Timeout)
	defer cancel()

	go func() {
		result <- g.Generate(options, &ctx)
	}()

	// Wait for the result or timeout
	err := <-result

	return err
}

func (g *Generator) Generate(options DebrickedOptions, ctx cgexec.IContext) error {
	targetPath := ".debrickedTmpFolder"
	debrickedExclusions := []string{targetPath}
	exclusions := append(options.Exclusions, debrickedExclusions...)

	var jobs []job.IJob
	for _, config := range options.Configs {
		s, strategyErr := g.strategyFactory.Make(config, options.Paths, exclusions, options.Inclusions, ctx)
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
