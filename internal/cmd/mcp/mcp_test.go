package mcp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMCPCmd(t *testing.T) {
	token := ""
	cmd := NewMCPCmd(&token, fakeAuthenticator{}, "https://debricked.com")
	commands := cmd.Commands()
	nbrOfCommands := 4
	assert.Lenf(t, commands, nbrOfCommands, "failed to assert that there were %d sub commands connected", nbrOfCommands)

	expected := map[string]bool{"start": false, "setup": false, "status": false, "uninstall": false}
	for _, subCmd := range commands {
		if _, ok := expected[subCmd.Name()]; ok {
			expected[subCmd.Name()] = true
		}
	}
	for name, found := range expected {
		assert.Truef(t, found, "failed to assert that %s command was registered", name)
	}
}
