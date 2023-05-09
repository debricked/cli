package callgraph

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/debricked/cli/pkg/callgraph/config"
	"github.com/debricked/cli/pkg/callgraph/job"
	"github.com/debricked/cli/pkg/callgraph/strategy"
	"github.com/debricked/cli/pkg/io/finder"
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

func findFiles(roots []string, exclusions []string) ([]string, error) {
	files := make(map[string]bool)
	var err error = nil

	for _, root := range roots {
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			for _, dir := range exclusions {
				if info.IsDir() && info.Name() == dir {
					return filepath.SkipDir
				}
			}

			if !info.IsDir() {
				files[path] = true
			}

			return nil
		})

		if err != nil {
			break
		}
	}

	fileList := make([]string, len(files))
	i := 0
	for k := range files {
		fileList[i] = k
		i++
	}

	return fileList, err
}

func (r Generator) Generate(paths []string, exclusions []string, configs []config.IConfig, status chan bool) (IGeneration, error) {
	targetPath := ".debrickedTmpFolder"
	debrickedExclusions := []string{targetPath}
	exclusions = append(exclusions, debrickedExclusions...)
	files, err := finder.FindFiles(paths, exclusions)
	fmt.Println(err)

	var jobs []job.IJob
	for _, config := range configs {
		s, strategyErr := r.strategyFactory.Make(config, files)
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
