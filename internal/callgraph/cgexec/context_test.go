package cgexec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDone(t *testing.T) {

	ctx, cancel := NewContext(0)
	defer cancel()
	val := <-ctx.Done()
	assert.NotNil(t, val)
}
