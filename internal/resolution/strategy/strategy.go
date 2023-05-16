package strategy

import (
	"github.com/debricked/cli/internal/resolution/job"
)

type IStrategy interface {
	Invoke() ([]job.IJob, error)
}
