package strategy

import (
	"github.com/debricked/cli/internal/callgraph/job"
)

type IStrategy interface {
	Invoke() ([]job.IJob, error)
}
