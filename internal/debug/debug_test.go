package debug

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	rescueStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Log("hello", true)

	_ = w.Close()
	output, _ := io.ReadAll(r)
	os.Stderr = rescueStderr

	assert.Contains(t, string(output), "DEBUG: ")
	assert.Contains(t, string(output), "hello\n")
}

func TestLogDebugDisabled(t *testing.T) {
	rescueStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Log("hello", false)

	_ = w.Close()
	output, _ := io.ReadAll(r)
	os.Stderr = rescueStderr

	assert.Empty(t, string(output))
}
