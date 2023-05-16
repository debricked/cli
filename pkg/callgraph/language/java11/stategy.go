package java

import (
	"fmt"
	"log"
	"path/filepath"

	conf "github.com/debricked/cli/pkg/callgraph/config"
	"github.com/debricked/cli/pkg/callgraph/job"
	"github.com/debricked/cli/pkg/io/finder"
	"github.com/debricked/cli/pkg/io/writer"
	"github.com/fatih/color"
)

type Strategy struct {
	config conf.IConfig
	files  []string
	finder finder.IFinder
}

func (s Strategy) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob
	// Filter relevant files

	if s.config == nil {
		strategyWarning("No config is setup")
		return jobs, nil
	}

	pmConfig := s.config.Kwargs()["pm"]

	var roots []string
	var err error
	switch pmConfig {
	case maven:
		roots, err = s.finder.FindMavenRoots(s.files)
	default:
		roots, err = s.finder.FindMavenRoots(s.files)
	}

	if err != nil {
		strategyWarning("Error while finding roots: " + err.Error())
		return jobs, nil
	}

	// TODO: If we want to build, build jobs need to execute before trying to find javaClassDirs.
	// If not, mapping between roots and classes could get wonky
	// Perfect time to build after getting roots, and maybe if no classes are found?

	classDirs, _ := s.finder.FindJavaClassDirs(s.files)
	absRoots, _ := finder.ConvertPathsToAbsPaths(roots)
	absClassDirs, _ := finder.ConvertPathsToAbsPaths(classDirs)
	rootClassMapping := finder.MapFilesToDir(absRoots, absClassDirs)

	for _, root := range absRoots {
		if _, ok := rootClassMapping[root]; ok == false {
			strategyWarning("Root found without related classes, make sure to build your project before running, root: " + root)
		}
	}
	if len(rootClassMapping) == 0 {
		return jobs, nil
	}

	for rootFile, classDirs := range rootClassMapping {
		// For each class paths dir within the root, find GCDPath as entrypoint
		classDir := finder.GCDPath(classDirs)
		rootDir := filepath.Dir(rootFile)
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

func strategyWarning(errMsg string) {
	err := fmt.Errorf(errMsg)
	warningColor := color.New(color.FgYellow, color.Bold).SprintFunc()
	defaultOutputWriter := log.Writer()
	log.Println(warningColor("Warning: ") + err.Error())
	log.SetOutput(defaultOutputWriter)
}
