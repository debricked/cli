package callgraph

import (
	"fmt"
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

const callgraph = "Callgraph"

func NewScheduler(workers int) *Scheduler {
	return &Scheduler{workers: workers, waitGroup: sync.WaitGroup{}}
}

func (scheduler *Scheduler) Schedule(jobs []job.IJob) (IGeneration, error) {
	fmt.Println("Starting scheduler")
	scheduler.queue = make(chan queueItem, len(jobs))
	scheduler.waitGroup.Add(len(jobs))
	scheduler.spinnerManager = tui.NewSpinnerManager()
	scheduler.spinnerManager.Start()
	fmt.Println("Done with spinner start", len(jobs))

	for _, j := range jobs {
		spinner := scheduler.spinnerManager.AddSpinner(callgraph, j.GetDir())
		scheduler.queue <- queueItem{
			job:     j,
			spinner: spinner,
		}
	}

	// Run it in sequence
	for item := range scheduler.queue {
		fmt.Println("start job")
		go scheduler.updateStatus(item)
		item.job.Run()
		scheduler.finish(item)
		scheduler.waitGroup.Done()
	}

	scheduler.spinnerManager.Stop()
	fmt.Println("Done")

	return NewGeneration(jobs), nil
}

func (scheduler *Scheduler) updateStatus(item queueItem) {
	for {
		msg := <-item.job.ReceiveStatus()
		tui.SetSpinnerMessage(item.spinner, callgraph, item.job.GetDir(), msg)
	}
}

func (scheduler *Scheduler) finish(item queueItem) {
	if item.job.Errors().HasError() {
		tui.SetSpinnerMessage(item.spinner, callgraph, item.job.GetDir(), "failed")
		item.spinner.Error()
	} else {
		tui.SetSpinnerMessage(item.spinner, callgraph, item.job.GetDir(), "done")
		item.spinner.Complete()
	}
}
