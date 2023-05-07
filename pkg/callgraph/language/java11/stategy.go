package java

import (
	"path/filepath"

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
	pattern := "*.jar"
	dirsWithJarFiles := make(map[string]bool)
	for _, file := range s.files {
		matched, _ := filepath.Match(pattern, filepath.Base(file))
		if matched {
			dirsWithJarFiles[filepath.Dir(file)] = true
		}
	}

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
