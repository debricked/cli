package tui

import (
	"fmt"
	"sync/atomic"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// ProgressBar tracks scan progress. It closes out with a checkmark when the
// scan it tracks completed, and with an error mark when it did not - a scan
// that is still queued once polling gives up must not read as a success.
//
// The embedded bar is safe for concurrent use, so failed is atomic to keep the
// wrapper safe too: the completion callback reads it on whichever goroutine
// happened to complete the bar.
type ProgressBar struct {
	*progressbar.ProgressBar
	failed atomic.Bool
}

func NewProgressBar() *ProgressBar {
	bar := &ProgressBar{}
	bar.ProgressBar = progressbar.NewOptions(100,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription("[blue]Scanning...[reset]"),
		progressbar.OptionOnCompletion(func() {
			color.NoColor = false
			if bar.failed.Load() {
				fmt.Println(color.RedString("⨯"))

				return
			}
			checkmark := color.GreenString("✔")
			fmt.Println(checkmark)
		}),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[blue]█[reset]",
			SaucerPadding: " ",
			BarStart:      "|",
			BarEnd:        "|",
		}),
	)

	return bar
}

// Fail leaves the bar at the progress it actually reached and closes it with an
// error mark, rather than filling it to 100% the way Finish does.
//
// Unlike Finish it does not advance the bar to its maximum, so IsFinished stays
// false and a `for !bar.IsFinished()` loop will not terminate on it. Callers
// must break out of such a loop themselves. The bar is spent afterwards -
// further Set/Add calls are no-ops.
func (p *ProgressBar) Fail() error {
	p.failed.Store(true)

	return p.Exit()
}
