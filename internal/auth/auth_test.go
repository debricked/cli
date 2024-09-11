package auth

import (
	"testing"

	"github.com/debricked/cli/internal/client/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewGeneration(t *testing.T) {
	res := NewDebrickedAuthenticator(testdata.NewDebClientMock())
	assert.NotNil(t, res)
}
