package tui

import (
	"fmt"
	"github.com/chelnak/ysmrr"
	"github.com/chelnak/ysmrr/pkg/colors"
	"github.com/fatih/color"
)

func NewSpinnerManager() ysmrr.SpinnerManager {
	return ysmrr.NewSpinnerManager(ysmrr.WithSpinnerColor(colors.FgHiBlue))
}

func SetSpinnerMessage(spinner *ysmrr.Spinner, filename string, message string) {
	spinner.UpdateMessage(fmt.Sprintf("Resolving %s: %s", color.YellowString(filename), message))
}
