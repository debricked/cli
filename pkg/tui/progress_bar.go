package tui

import (
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

func NewProgressBar() *progressbar.ProgressBar {
	return progressbar.NewOptions(100,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionSetWidth(30),
		progressbar.OptionSetDescription("[blue]Scanning...[reset]"),
		progressbar.OptionOnCompletion(func() {
			color.NoColor = false
			color.Green("✔")
		}),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[blue]█[reset]",
			SaucerPadding: " ",
			BarStart:      "|",
			BarEnd:        "|",
		}),
	)
}
