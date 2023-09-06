package nuget

import (
	"github.com/debricked/cli/internal/resolution/job"
)

type Strategy struct {
	files []string
}

func (s Strategy) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob
	for _, file := range s.files {
		jobs = append(jobs, NewJob(
			file,
			true,
			CmdFactory{
				execPath: ExecPath{},
			},
		),
		)
	}

	return jobs, nil
}

func NewStrategy(files []string) Strategy {
	return Strategy{files}
}
