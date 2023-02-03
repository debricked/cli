package strategy

import (
	"github.com/debricked/cli/pkg/resolution/job"
)

type IStrategy interface {
	Invoke() []job.IJob
}
