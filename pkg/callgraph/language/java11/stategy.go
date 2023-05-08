package java

import (
	conf "github.com/debricked/cli/pkg/callgraph/config"
	"github.com/debricked/cli/pkg/callgraph/job"
	"github.com/debricked/cli/pkg/io/finder"
	"github.com/debricked/cli/pkg/io/writer"
)

type Strategy struct {
	config conf.IConfig
	files  []string
}

func (s Strategy) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob
	// Filter relevant files

	roots := finder.FindMavenRoots(s.files)
	jarFiles := finder.FindJarDirs(s.files)
	rootPomJarMapping := finder.MapFilesToDir(roots, jarFiles)

	for rootDir, jars := range rootPomJarMapping {
		jobs = append(jobs, NewJob(
			rootDir,
			jars,
			CmdFactory{},
			writer.FileWriter{},
			s.config,
		),
		)
	}

	return jobs, nil
}

func NewStrategy(config conf.IConfig, files []string) Strategy {
	return Strategy{config, files}
}
