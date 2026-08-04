package tui

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewProgressBar(t *testing.T) {
	bar := NewProgressBar()

	assert.NotNil(t, bar)
	assert.False(t, bar.IsFinished(), "failed to assert that the bar was not finished")

	err := bar.Set(100)
	assert.NoError(t, err)
	assert.True(t, bar.IsFinished(), "failed to assert that the bar was finished")
}

func TestProgressBarFail(t *testing.T) {
	bar := NewProgressBar()
	err := bar.Set(40)
	assert.NoError(t, err)

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = bar.Fail()

	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	assert.NoError(t, err)
	assert.Contains(t, string(out), "⨯", "failed to assert that a failed bar renders an error mark")
	assert.NotContains(t, string(out), "✔", "failed to assert that a failed bar renders no checkmark")
	assert.False(t, bar.IsFinished(), "failed to assert that a failed bar was not filled to 100%")
}
