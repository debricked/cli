package tui

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestNewSpinnerManager(t *testing.T) {
	spinnerManager := NewSpinnerManager(
		"Resolving",
		"waiting for worker",
	)

	assert.NotNil(t, spinnerManager)
}

func TestSetSpinnerMessage(t *testing.T) {
	spinnerManager := NewSpinnerManager(
		"Resolving",
		"waiting for worker",
	)
	message := "test"
	spinner := spinnerManager.AddSpinner(message)
	assert.Contains(t, spinner.GetMessage(), fmt.Sprintf("Resolving %s: waiting for worker", color.YellowString(message)))

	fileName := "file-name"
	message = "new test message"

	spinnerManager.SetSpinnerMessage(spinner, fileName, message)
	assert.Contains(t, spinner.GetMessage(), fmt.Sprintf("Resolving %s: %s", color.YellowString(fileName), message))
}

func TestSetSpinnerMessageLongFilenameParts(t *testing.T) {
	spinnerManager := NewSpinnerManager(
		"Resolving",
		"waiting for worker",
	)
	longFilenameParts := []string{
		"directory",
		"sub-directory################################################################",
		"file.json",
	}
	longFileName := filepath.Join(longFilenameParts...)

	spinner := spinnerManager.AddSpinner(longFileName)
	message := spinner.GetMessage()

	assert.Contains(t, message, longFileName)
}

func TestSetSpinnerMessageLongFilenameManyDirs(t *testing.T) {
	spinnerManager := NewSpinnerManager(
		"Resolving",
		"waiting for worker",
	)
	longFilenameParts := []string{
		"directory",
		"sub-directory",
		"sub-directory",
		"sub-directory",
		"sub-directory",
		"sub-directory",
		"target-directory",
		"file.json",
	}
	longFileName := filepath.Join(longFilenameParts...)

	truncatedFilenameParts := []string{
		longFilenameParts[0],
		"...",
		longFilenameParts[len(longFilenameParts)-2],
		longFilenameParts[len(longFilenameParts)-1],
	}
	truncatedFilename := filepath.Join(truncatedFilenameParts...)
	spinner := spinnerManager.AddSpinner(longFileName)
	message := spinner.GetMessage()

	assert.Contains(t, message, truncatedFilename)
}

func TestStartStop(t *testing.T) {
	spinnerManager := NewSpinnerManager(
		"Resolving",
		"waiting for worker",
	)

	spinnerManager.Start()
	spinnerManager.Stop()
}
