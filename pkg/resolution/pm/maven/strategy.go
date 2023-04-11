package maven

import (
	"github.com/debricked/cli/pkg/resolution/job"
)

type Strategy struct {
	files      []string
	cmdFactory ICmdFactory
	pomService IPomService
}

func NewStrategy(files []string) Strategy {
	return Strategy{files, CmdFactory{}, PomService{}}
}

func (s Strategy) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob
	s.files = s.pomService.GetRootPomFiles(s.files)

	for _, file := range s.files {
		jobs = append(jobs, NewJob(file, s.cmdFactory))
	}

	return jobs, nil
}
