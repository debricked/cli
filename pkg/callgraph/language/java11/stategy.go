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
	finder finder.IFinder
}

func (s Strategy) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob
	// Filter relevant files

	pmConfig := s.config.Kwargs()["pm"]
	println("CONFIG", pmConfig)
	var roots []string
	var err error
	switch pmConfig {
	case gradle:
		roots, err = s.finder.FindGradleRoots(s.files)
	case maven:
		roots, err = s.finder.FindMavenRoots(s.files)
	default:
		roots, err = s.finder.FindMavenRoots(s.files)
	}

	if err != nil {
		fmt.Println("error", err)
		return jobs, err
	}

	// TODO: If we want to build, build jobs need to execute before trying to find javaClassDirs.
	// If not, mapping between roots and classes could get wonky
	// Perfect time to build after getting roots, and maybe if no classes are found?

	fmt.Println("roots", roots)
	classDirs, _ := s.finder.FindJavaClassDirs(s.files)
	rootClassMapping := s.finder.MapFilesToDir(roots, classDirs)
	fmt.Println("roots", rootClassMapping)

	if len(roots) != 0 && len(rootClassMapping) == 0 {
		fmt.Println("error", err)
		return jobs, fmt.Errorf("Roots found but without related classes, make sure to build your project before running")
	}

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

func NewStrategy(config conf.IConfig, files []string, finder finder.IFinder) Strategy {
	return Strategy{config, files, finder}
}
