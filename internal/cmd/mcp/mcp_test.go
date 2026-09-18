package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMCPCmd(t *testing.T) {
	token := ""
	cmd := NewMCPCmd(&token, fakeAuthenticator{}, "https://debricked.com")
	commands := cmd.Commands()
	nbrOfCommands := 1
	assert.Lenf(t, commands, nbrOfCommands, "failed to assert that there were %d sub commands connected", nbrOfCommands)

	match := false
	for _, subCmd := range commands {
		if subCmd.Name() == "start" {
			match = true

			break
		}
	}
	assert.True(t, match, "failed to assert that start command was registered")
}
