package tui

import (
	"fmt"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func TestNewSpinnerManager(t *testing.T) {
	spinnerManager := NewSpinnerManager()
	assert.NotNil(t, spinnerManager)
}

func TestSetSpinnerMessage(t *testing.T) {
	spinnerManager := NewSpinnerManager()
	message := "test"
	spinner := spinnerManager.AddSpinner(message)
	assert.Contains(t, spinner.GetMessage(), fmt.Sprintf("Resolving %s: waiting for worker", color.YellowString(message)))

	fileName := "file-name"
	message = "new test message"

	SetSpinnerMessage(spinner, fileName, message)
	assert.Contains(t, spinner.GetMessage(), fmt.Sprintf("Resolving %s: %s", color.YellowString(fileName), message))
}

func TestStartStop(t *testing.T) {
	spinnerManager := NewSpinnerManager()
	spinnerManager.Start()
	spinnerManager.Stop()
}
