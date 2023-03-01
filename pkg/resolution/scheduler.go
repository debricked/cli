package resolution

import (
	"sync"

	"github.com/chelnak/ysmrr"
	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/tui"
)

type IScheduler interface {
	Schedule(jobs []job.IJob) (IResolution, error)
}

type queueItem struct {
	job     job.IJob
	spinner *ysmrr.Spinner
}

type Scheduler struct {
	workers        int
	queue          chan queueItem
	waitGroup      sync.WaitGroup
	spinnerManager ysmrr.SpinnerManager
}

func NewScheduler(workers int) *Scheduler {
	return &Scheduler{workers: workers, waitGroup: sync.WaitGroup{}}
}

func (scheduler *Scheduler) Schedule(jobs []job.IJob) (IResolution, error) {
	scheduler.queue = make(chan queueItem, len(jobs))
	scheduler.waitGroup.Add(len(jobs))

	scheduler.spinnerManager = tui.NewSpinnerManager()
	for range jobs {
		scheduler.spinnerManager.AddSpinner("")
	}
	scheduler.spinnerManager.Start()

	for w := 1; w <= scheduler.workers; w++ {
		go scheduler.worker()
	}

	for i := range jobs {
		scheduler.queue <- queueItem{
			job:     jobs[i],
			spinner: scheduler.spinnerManager.GetSpinners()[i],
		}
	}

	scheduler.waitGroup.Wait()

	scheduler.spinnerManager.Stop()

	close(scheduler.queue)

	return NewResolution(jobs), nil
}

func (scheduler *Scheduler) worker() {
	for item := range scheduler.queue {
		go scheduler.updateStatus(item)

		item.job.Run()

		scheduler.finish(item)

		scheduler.waitGroup.Done()
	}
}
func (scheduler *Scheduler) updateStatus(item queueItem) {
	for {
		msg := <-item.job.ReceiveStatus()
		tui.SetSpinnerMessage(item.spinner, item.job.GetFile(), msg)
	}
}

func (scheduler *Scheduler) finish(item queueItem) {
	if item.job.Error() == nil {
		tui.SetSpinnerMessage(item.spinner, item.job.GetFile(), "done")
		item.spinner.Complete()
	} else {
		tui.SetSpinnerMessage(item.spinner, item.job.GetFile(), "failed")
		item.spinner.Error()
	}
}
