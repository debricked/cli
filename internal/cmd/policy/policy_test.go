package policy

import (
	"testing"

	"github.com/debricked/cli/internal/cmd/policy/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewPolicyCmd(t *testing.T) {
	cmd := NewPolicyCmd(&testdata.ValidatorMock{})
	commands := cmd.Commands()

	assert.Len(t, commands, 1)
	assert.Equal(t, "validate", commands[0].Name())
}

func TestPreRun(t *testing.T) {
	cmd := NewPolicyCmd(&testdata.ValidatorMock{})
	cmd.PreRun(cmd, nil)
}
