package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chelnak/ysmrr"
	"github.com/chelnak/ysmrr/pkg/colors"
	"github.com/fatih/color"
)

type ISpinnerManager interface {
	AddSpinner(action string, file string) *ysmrr.Spinner
	Start()
	Stop()
}

type SpinnerManager struct {
	spinnerManager ysmrr.SpinnerManager
}

func NewSpinnerManager() SpinnerManager {
	return SpinnerManager{ysmrr.NewSpinnerManager(ysmrr.WithSpinnerColor(colors.FgHiBlue))}
}

func (sm SpinnerManager) AddSpinner(action string, file string) *ysmrr.Spinner {
	spinner := sm.spinnerManager.AddSpinner("")
	SetSpinnerMessage(spinner, action, file, "waiting for worker")

	return spinner
}

func (sm SpinnerManager) Start() {
	sm.spinnerManager.Start()
}

func (sm SpinnerManager) Stop() {
	sm.spinnerManager.Stop()
}

func SetSpinnerMessage(spinner *ysmrr.Spinner, action string, filename string, message string) {
	const maxNumberOfChars = 50
	truncatedFilename := filename
	if len(truncatedFilename) > maxNumberOfChars {
		separator := string(os.PathSeparator)
		pathParts := strings.Split(filename, separator)
		if len(pathParts) > 3 {
			firstDir := pathParts[0]
			lastDir := pathParts[len(pathParts)-2]
			name := pathParts[len(pathParts)-1]
			truncatedFilename = filepath.Join(
				firstDir,
				"...",
				lastDir,
				name,
			)
		}

	}
	file := color.YellowString(truncatedFilename)
	spinner.UpdateMessage(fmt.Sprintf("%s %s: %s", action, file, message))
}
