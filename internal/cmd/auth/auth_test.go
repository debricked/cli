package auth

import (
	"testing"

	"github.com/debricked/cli/internal/auth"
	"github.com/debricked/cli/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestNewFilesCmd(t *testing.T) {
	token := "token"
	deb_client := client.NewDebClient(&token, nil)
	authenticator := auth.NewDebrickedAuthenticator(deb_client)
	cmd := NewAuthCmd(authenticator)
	commands := cmd.Commands()
	nbrOfCommands := 3
	assert.Lenf(t, commands, nbrOfCommands, "failed to assert that there were %d sub commands connected", nbrOfCommands)
}

func TestPreRun(t *testing.T) {
	cmd := NewAuthCmd(nil)
	cmd.PreRun(cmd, nil)
}
