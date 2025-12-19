package token

import (
	"testing"

	"github.com/debricked/cli/internal/auth"
	"github.com/debricked/cli/internal/auth/testdata"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewTokenCmd(t *testing.T) {
	authenticator := auth.NewDebrickedAuthenticator("")
	cmd := NewTokenCmd(authenticator)
	commands := cmd.Commands()
	nbrOfCommands := 0
	assert.Len(t, commands, nbrOfCommands)

	flags := cmd.Flags()
	flagAssertions := map[string]string{}
	for name, shorthand := range flagAssertions {
		flag := flags.Lookup(name)
		assert.NotNil(t, flag)
		assert.Equal(t, shorthand, flag.Shorthand)
	}

	var flagKeys = []string{
		JsonFlag,
	}
	viperKeys := viper.AllKeys()
	for _, flagKey := range flagKeys {
		match := false
		for _, key := range viperKeys {
			if key == flagKey {
				match = true
			}
		}
		assert.Truef(t, match, "failed to assert that flag was present: "+flagKey)
	}
}

func TestPreRun(t *testing.T) {
	mockAuthenticator := testdata.MockAuthenticator{}
	cmd := NewTokenCmd(mockAuthenticator)
	cmd.PreRun(cmd, nil)
}

func TestRunE(t *testing.T) {
	a := testdata.MockAuthenticator{}
	runE := RunE(a)

	err := runE(nil, []string{})

	assert.NoError(t, err)
}

func TestRunEJSONFlag(t *testing.T) {
	a := testdata.MockAuthenticator{}
	jsonFormat = true
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
