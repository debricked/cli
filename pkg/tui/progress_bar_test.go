package tui

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewProgressBar(t *testing.T) {
	bar := NewProgressBar()

	assert.NotNil(t, bar)
	assert.False(t, bar.IsFinished(), "failed to assert that the bar was not finished")

	err := bar.Set(100)
	assert.NoError(t, err)
	assert.True(t, bar.IsFinished(), "failed to assert that the bar was finished")
}
