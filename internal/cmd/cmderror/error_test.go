package cmderror

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := CommandError{Code: 2, Err: errors.New("failure")}

	assert.Equal(t, "failure", err.Error())
}

func TestUnwrap(t *testing.T) {
	wrappedErr := errors.New("failure")
	err := CommandError{Code: 2, Err: fmt.Errorf("context. %w", wrappedErr)}

	assert.ErrorIs(t, err, wrappedErr)

	var cmdErr CommandError
	assert.ErrorAs(t, error(err), &cmdErr)
	assert.Equal(t, 2, cmdErr.Code)
}
