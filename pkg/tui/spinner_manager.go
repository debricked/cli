package tui

import (
	"fmt"

	"github.com/chelnak/ysmrr"
	"github.com/chelnak/ysmrr/pkg/colors"
	"github.com/fatih/color"
)

type ISpinnerManager interface {
	AddSpinner(file string) *ysmrr.Spinner
	Start()
	Stop()
}

type SpinnerManager struct {
	spinnerManager ysmrr.SpinnerManager
}

func NewSpinnerManager() SpinnerManager {
	return SpinnerManager{ysmrr.NewSpinnerManager(ysmrr.WithSpinnerColor(colors.FgHiBlue))}
}

func (sm SpinnerManager) AddSpinner(file string) *ysmrr.Spinner {
	spinner := sm.spinnerManager.AddSpinner("")
	SetSpinnerMessage(spinner, file, "waiting for worker")

	return spinner
}

func (sm SpinnerManager) Start() {
	sm.spinnerManager.Start()
}

func (sm SpinnerManager) Stop() {
	sm.spinnerManager.Stop()
}

func SetSpinnerMessage(spinner *ysmrr.Spinner, filename string, message string) {
	file := color.YellowString(filename)
	spinner.UpdateMessage(fmt.Sprintf("Resolving %s: %s", file, message))
}
