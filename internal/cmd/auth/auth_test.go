package auth

import (
	"testing"

	"github.com/debricked/cli/internal/auth"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthCmd(t *testing.T) {
	authenticator := auth.NewDebrickedAuthenticator("")
	cmd := NewAuthCmd(authenticator)
	commands := cmd.Commands()
	nbrOfCommands := 3
	assert.Lenf(t, commands, nbrOfCommands, "failed to assert that there were %d sub commands connected", nbrOfCommands)
}

func TestPreRun(t *testing.T) {
	cmd := NewAuthCmd(nil)
	cmd.PreRun(cmd, nil)
}
