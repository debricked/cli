package maven

import (
	"github.com/debricked/cli/pkg/resolution/job"
)

type Strategy struct {
	files      []string
	cmdFactory ICmdFactory
	pomX       IPomX
}

func NewStrategy(files []string) Strategy {

	return Strategy{files, CmdFactory{}, PomX{}}
}

func (s Strategy) Invoke() []job.IJob {

	var jobs []job.IJob

	// filter out the root pom files
	s.files = s.pomX.GetRootPomFiles(s.files)

	for _, file := range s.files {
		jobs = append(jobs, NewJob(file, s.cmdFactory))
	}

	return jobs
}
