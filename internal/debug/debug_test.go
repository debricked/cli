package debug

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrint(t *testing.T) {
	rescueStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Print("hello", true)

	_ = w.Close()
	output, _ := io.ReadAll(r)
	os.Stderr = rescueStderr

	assert.Equal(t, "DEBUG: hello\n", string(output))
}
