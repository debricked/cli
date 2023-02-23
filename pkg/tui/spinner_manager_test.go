package tui

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewSpinnerManager(t *testing.T) {
	spinnerManager := NewSpinnerManager()
	assert.NotNil(t, spinnerManager)
}

func TestSetSpinnerMessage(t *testing.T) {
	spinnerManager := NewSpinnerManager()
	message := "test"
	spinner := spinnerManager.AddSpinner(message)
	assert.Equal(t, message, spinner.GetMessage())

	fileName := "file-name"
	message = "new test message"

	SetSpinnerMessage(spinner, fileName, message)
	assert.Contains(t, spinner.GetMessage(), fmt.Sprintf("Resolving %s: %s", color.YellowString(fileName), message))
}
