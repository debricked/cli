package callgraph

import (
	"sync"

	"github.com/chelnak/ysmrr"
	"github.com/debricked/cli/pkg/callgraph/job"
	"github.com/debricked/cli/pkg/tui"
)

type IScheduler interface {
	Schedule(jobs []job.IJob) (IGeneration, error)
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

func (scheduler *Scheduler) Schedule(jobs []job.IJob) (IGeneration, error) {
	scheduler.queue = make(chan queueItem, len(jobs))
	scheduler.waitGroup.Add(len(jobs))
	scheduler.spinnerManager = tui.NewSpinnerManager()
	scheduler.spinnerManager.Start()

	// Run it in sequence
	for item := range scheduler.queue {
		go scheduler.updateStatus(item)
		item.job.Run()
		scheduler.finish(item)
		scheduler.waitGroup.Done()
	}

	scheduler.spinnerManager.Stop()

	return NewGeneration(jobs), nil
}

func (scheduler *Scheduler) updateStatus(item queueItem) {
	for {
		msg := <-item.job.ReceiveStatus()
		tui.SetSpinnerMessage(item.spinner, "file", msg)
	}
}

func (scheduler *Scheduler) finish(item queueItem) {
	if item.job.Errors().HasError() {
		tui.SetSpinnerMessage(item.spinner, "file", "failed")
		item.spinner.Error()
	} else {
		tui.SetSpinnerMessage(item.spinner, "file", "done")
		item.spinner.Complete()
	}
}
