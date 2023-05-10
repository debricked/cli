package resolution

import (
	"sort"
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
	spinnerManager tui.ISpinnerManager
}

const resolving = "Resolving"

func NewScheduler(workers int) *Scheduler {
	return &Scheduler{workers: workers, waitGroup: sync.WaitGroup{}}
}

func (scheduler *Scheduler) Schedule(jobs []job.IJob) (IResolution, error) {
	scheduler.queue = make(chan queueItem, len(jobs))
	scheduler.waitGroup.Add(len(jobs))

	scheduler.spinnerManager = tui.NewSpinnerManager()

	for w := 1; w <= scheduler.workers; w++ {
		go scheduler.worker()
	}

	sort.Slice(jobs, func(i, j int) bool {
		return jobs[i].GetFile() < jobs[j].GetFile()
	})

	for _, j := range jobs {
		spinner := scheduler.spinnerManager.AddSpinner(resolving, j.GetFile())
		scheduler.queue <- queueItem{
			job:     j,
			spinner: spinner,
		}
	}
	scheduler.spinnerManager.Start()

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
		tui.SetSpinnerMessage(item.spinner, resolving, item.job.GetFile(), msg)
	}
}

func (scheduler *Scheduler) finish(item queueItem) {
	if item.job.Errors().HasError() {
		tui.SetSpinnerMessage(item.spinner, resolving, item.job.GetFile(), "failed")
		item.spinner.Error()
	} else {
		tui.SetSpinnerMessage(item.spinner, resolving, item.job.GetFile(), "done")
		item.spinner.Complete()
	}
}
