package maven

import (
	"path/filepath"

	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/vifraa/gopom"
)

type Strategy struct {
	files      []string
	cmdFactory ICmdFactory
}

func NewStrategy(files []string) Strategy {
	return Strategy{files, CmdFactory{}}
}

func (s Strategy) Invoke() []job.IJob {

	// todo make use of the gopom as an interface

	var jobs []job.IJob

	rootPaths := s.GetRootPomFiles()

	for _, file := range rootPaths {
		jobs = append(jobs, NewJob(file, s.cmdFactory))
	}

	return jobs
}

func (s Strategy) GetRootPomFiles() []string {

	childMap := make(map[string][]string)

	for _, file_path := range s.files {

		pom, _ := gopom.Parse(file_path)

		if len(pom.Modules) == 0 {
			continue
		}

		for _, module := range pom.Modules {

			// path to child pom
			modulePath := filepath.Join(filepath.Dir(file_path), filepath.Dir(module), filepath.Base(module), "pom.xml")

			childMap[modulePath] = append(childMap[modulePath], file_path)

		}
	}

	roots := make([]string, 0)

	for _, file := range s.files {
		if _, ok := childMap[file]; !ok {
			roots = append(roots, file)
		}
	}

	return roots

}
