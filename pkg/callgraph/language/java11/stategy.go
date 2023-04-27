package java

import (
	conf "github.com/debricked/cli/pkg/callgraph/config"
	"github.com/debricked/cli/pkg/callgraph/job"
	"github.com/debricked/cli/pkg/io/writer"
)

type Strategy struct {
	config conf.IConfig
	files  []string
}

func (s Strategy) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob
	// Filter relevant files
	jobs = append(jobs, NewJob(
		s.files,
		CmdFactory{},
		writer.FileWriter{},
		s.config,
	),
	)

	return jobs, nil
}

func NewStrategy(config conf.IConfig, files []string) Strategy {
	return Strategy{config, files}
}
