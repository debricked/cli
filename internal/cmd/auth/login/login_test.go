package login

import (
	"testing"

	"github.com/debricked/cli/internal/auth"
	"github.com/debricked/cli/internal/auth/testdata"
	"github.com/debricked/cli/internal/client"
	"github.com/stretchr/testify/assert"
)

func TestNewLoginCmd(t *testing.T) {
	token := "token"
	deb_client := client.NewDebClient(&token, nil)
	authenticator := auth.NewDebrickedAuthenticator(deb_client)
	cmd := NewLoginCmd(authenticator)
	commands := cmd.Commands()
	nbrOfCommands := 0
	assert.Len(t, commands, nbrOfCommands)
}

func TestPreRun(t *testing.T) {
	mockAuthenticator := testdata.MockAuthenticator{}
	cmd := NewLoginCmd(mockAuthenticator)
	cmd.PreRun(cmd, nil)
}

func TestRunE(t *testing.T) {
	a := testdata.MockAuthenticator{}
	runE := RunE(a)

	err := runE(nil, []string{})

	assert.NoError(t, err)
}

func TestRunEError(t *testing.T) {
	a := testdata.ErrorMockAuthenticator{}
	runE := RunE(a)

	err := runE(nil, []string{})

	assert.Error(t, err)
}
