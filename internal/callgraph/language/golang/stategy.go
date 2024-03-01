package golang

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	conf "github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/callgraph/finder"
	"github.com/debricked/cli/internal/callgraph/job"
	"github.com/debricked/cli/internal/io"
	"github.com/fatih/color"
)

type Strategy struct {
	config     conf.IConfig
	cmdFactory ICmdFactory
	paths      []string
	exclusions []string
	finder     finder.IFinder
	ctx        cgexec.IContext
}

func (s Strategy) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob

	if s.config == nil {
		strategyWarning("No config is setup")

		return jobs, nil
	}

	for _, path := range s.paths {

		files, err := s.finder.FindFiles([]string{path}, s.exclusions)
		if err != nil {
			strategyWarning("Error while finding files: " + err.Error())

			return jobs, err
		}

		roots, err := s.finder.FindRoots(files)
		if err != nil {
			strategyWarning("Error while finding roots: " + err.Error())

			return jobs, err
		}

		if len(roots) == 0 {
			strategyWarning("No main.go found")

			return jobs, nil
		}

		for _, rootFilePath := range roots {

			rootFileDir := filepath.Dir(rootFilePath)
			rootFile := filepath.Base(rootFilePath)

			jobs = append(jobs, NewJob(
				rootFileDir,
				rootFile,
				s.cmdFactory,
				io.FileWriter{},
				io.NewArchive("."),
				s.config,
				s.ctx,
				io.FileSystem{},
			),
			)
		}
	}
	return jobs, nil
}

func NewStrategy(config conf.IConfig, paths []string, exclusions []string, finder finder.IFinder, ctx cgexec.IContext) Strategy {
	return Strategy{config, CmdFactory{}, paths, exclusions, finder, ctx}
}

func strategyWarning(errMsg string) {
	err := fmt.Errorf(errMsg)
	warningColor := color.New(color.FgYellow, color.Bold).SprintFunc()
	defaultOutputWriter := log.Writer()
	log.Println(warningColor("Warning: ") + err.Error())
	log.SetOutput(defaultOutputWriter)
}
