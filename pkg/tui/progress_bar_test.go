package tui

import (
	"testing"
)

func TestNewProgressBar(t *testing.T) {
	bar := NewProgressBar()
	if bar == nil {
		t.Error("failed to assert bar was not nil")
	}

	if bar.IsFinished() {
		t.Error("failed to assert that the bar was not finished")
	}

	err := bar.Set(100)
	if err != nil {
		t.Error("failed to set progress to 100")
	}

	if !bar.IsFinished() {
		t.Error("failed to assert that the bar was finished")
	}
}
