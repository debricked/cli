package resolution

import (
	"github.com/debricked/cli/pkg/resolution/file"
	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/strategy"
)

type IResolver interface {
	Resolve(files []string) (IResolution, error)
}

type Resolver struct {
	batchFactory    file.IBatchFactory
	strategyFactory strategy.IFactory
	scheduler       IScheduler
}

func NewResolver(
	batchFactory file.IBatchFactory,
	strategyFactory strategy.IFactory,
	scheduler IScheduler,
) Resolver {
	return Resolver{
		batchFactory,
		strategyFactory,
		scheduler,
	}
}

func (r Resolver) Resolve(files []string) (IResolution, error) {
	pmBatches := r.batchFactory.Make(files)

	var jobs []job.IJob
	for _, pmBatch := range pmBatches {
		s, err := r.strategyFactory.Make(pmBatch)
		if err == nil {
			jobs = append(jobs, s.Invoke()...)
		}
	}

	resolution, err := r.scheduler.Schedule(jobs)

	return resolution, err
}
