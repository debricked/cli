package callgraph

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	"github.com/debricked/cli/pkg/callgraph/config"
	"github.com/debricked/cli/pkg/callgraph/job"
	"github.com/debricked/cli/pkg/callgraph/strategy"
)

type IGenerator interface {
	GenerateWithTimer(paths []string, exclusions []string, configs []config.IConfig, timeout int) error
	Generate(paths []string, exclusions []string, configs []config.IConfig, status chan bool) (IGeneration, error)
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

func (r Generator) GenerateWithTimer(paths []string, exclusions []string, configs []config.IConfig, timeout int) error {
	status := make(chan bool)
	timeoutChan := time.After(time.Duration(timeout) * time.Second)
	fmt.Println("Start generation")
	go r.Generate(paths, exclusions, configs, status)
	select {
	case <-status:
		fmt.Println("Function completed successfully")
	case <-timeoutChan:
		fmt.Println("Function timed out")
		// use the runtime package to kill the goroutine
		runtime.Goexit()
		return errors.New("Timeout reached, termingating generate callgraph goroutine")
	}

	return nil
}

func (r Generator) Generate(paths []string, exclusions []string, configs []config.IConfig, status chan bool) (IGeneration, error) {

	// For each config (or single root provided with commands from CMD), run CG.generation.

	// Might need to run it sequentially since it is a heavy operation that in itself might use multithreads

	// Find job-files,
	// Run scheduler on the jobs
	// add refine-path-step

	var jobs []job.IJob
	for _, config := range configs {
		fmt.Println("hello", config, paths)
		fmt.Println("strFac", r.strategyFactory)
		s, strategyErr := r.strategyFactory.Make(config, paths)
		if strategyErr == nil {
			newJobs, err := s.Invoke()
			if err != nil {
				return nil, err
			}
			jobs = append(jobs, newJobs...)
		}
	}

	fmt.Println("Run scheduler")
	generation, err := r.scheduler.Schedule(jobs)
	return generation, err
}
