package strategy

import (
	"github.com/debricked/cli/pkg/callgraph/job"
)

type IStrategy interface {
	Invoke() ([]job.IJob, error)
}
