package java

import (
	"fmt"

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

	pmConfig := s.config.Kwargs()["pm"]
	var roots []string
	var err error
	switch pmConfig {
	case gradle:
		roots, err = finder.FindGradleRoots(s.files)
	case maven:
		roots, err = finder.FindMavenRoots(s.files)
	default:
		roots, err = finder.FindMavenRoots(s.files)
	}

	if err != nil {
		fmt.Println("error", err)
	}

	classDirs := finder.FindJavaClassDirs(s.files)
	rootClassMapping := finder.MapFilesToDir(roots, classDirs)
	fmt.Println("roots", rootClassMapping)

	for rootDir, classDirs := range rootClassMapping {
		// For each class paths dir within the root, find GCDPath as entrypoint
		classDir := finder.GCDPath(classDirs)
		jobs = append(jobs, NewJob(
			rootDir,
			[]string{classDir},
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
