package callgraph

import (
	"sync"

	"github.com/chelnak/ysmrr"
	"github.com/debricked/cli/internal/callgraph/cgexec"
	"github.com/debricked/cli/internal/callgraph/job"
	"github.com/debricked/cli/internal/tui"
)

type IScheduler interface {
	Schedule(jobs []job.IJob, ctx cgexec.IContext) (IGeneration, error)
}

type queueItem struct {
	job     job.IJob
	spinner *ysmrr.Spinner
}

type Scheduler struct {
	workers        int
	queue          chan queueItem
	waitGroup      sync.WaitGroup
	spinnerManager tui.ISpinnerManager
}

func NewScheduler(workers int) *Scheduler {
	return &Scheduler{workers: workers, waitGroup: sync.WaitGroup{}}
}

func (scheduler *Scheduler) Schedule(jobs []job.IJob, ctx cgexec.IContext) (IGeneration, error) {
	if len(jobs) == 0 {
		return NewGeneration(jobs), nil
	}

	scheduler.queue = make(chan queueItem, len(jobs))
	scheduler.waitGroup.Add(len(jobs))
	scheduler.spinnerManager = tui.NewSpinnerManager("Calgraph", "waiting for worker")
	scheduler.spinnerManager.Start()

	for _, j := range jobs {
		spinner := scheduler.spinnerManager.AddSpinner(callgraph, j.GetDir())
		scheduler.queue <- queueItem{
			job:     j,
			spinner: spinner,
		}
	}

	jobIteration := 0
	// Run it in sequence
	for item := range scheduler.queue {
		jobIteration += 1
		go scheduler.updateStatus(item)
		item.job.Run()
		scheduler.finish(item)
		scheduler.waitGroup.Done()

		interupt := false
		if ctx != nil {
			select {
			case <-ctx.Done():
				close(scheduler.queue)
				interupt = true

				break
			default:
			}
		}

		if interupt {
			break
		}

		if jobIteration == len(jobs) {
			close(scheduler.queue)
		}
	}

	scheduler.spinnerManager.Stop()

	return NewGeneration(jobs), nil
}

func (scheduler *Scheduler) updateStatus(item queueItem) {
	for {
		msg := <-item.job.ReceiveStatus()
		scheduler.spinnerManager.SetSpinnerMessage(item.spinner, item.job.GetDir(), msg)
	}
}

func (scheduler *Scheduler) finish(item queueItem) {
	if item.job.Errors().HasError() {
		scheduler.spinnerManager.SetSpinnerMessage(item.spinner, item.job.GetDir(), "failed")
		item.spinner.Error()
	} else {
		scheduler.spinnerManager.SetSpinnerMessage(item.spinner, item.job.GetDir(), "done")
		item.spinner.Complete()
	}
}
