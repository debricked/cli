package callgraph

import (
	"github.com/debricked/cli/pkg/callgraph/config"
	"github.com/debricked/cli/pkg/callgraph/job"
	"github.com/debricked/cli/pkg/callgraph/strategy"
)

type IGenerator interface {
	Generate(paths []string, exclusions []string, configs []config.IConfig) (IGeneration, error)
}

type Generator struct {
	strategyFactory strategy.IFactory
	scheduler       IScheduler
}

func NewGenerator(
	strategyFactory strategy.IFactory,
	scheduler IScheduler,
) Generator {
	return Generator{
		strategyFactory,
		scheduler,
	}
}

func (r Generator) Generate(paths []string, exclusions []string, configs []config.IConfig) (IGeneration, error) {
	// Find roots we could potentialy care about

	// For each root (or single root provided with commands from CMD), run CG.generation.

	// Might need to run it sequentially since it is a heavy operation that in itself might use multithreads

	// Find job-files,
	// Run scheduler on the jobs

	var jobs []job.IJob
	for _, config := range configs {
		s, strategyErr := r.strategyFactory.Make(config, paths)
		if strategyErr == nil {
			newJobs, err := s.Invoke()
			if err != nil {
				return nil, err
			}
			jobs = append(jobs, newJobs...)
		}
	}

	resolution, err := r.scheduler.Schedule(jobs)

	// if resolution.HasErr() {
	// 	jobErrList := tui.NewJobsErrorList(os.Stdout, resolution.Jobs())
	// 	err = jobErrList.Render()
	// }

	return resolution, err
}
