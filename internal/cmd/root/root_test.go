package root

import (
	"bytes"
	"testing"

	"github.com/debricked/cli/internal/wire"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd("v0.0.0", wire.GetCliContainer())
	commands := cmd.Commands()
	nbrOfCommands := 7
	if len(commands) != nbrOfCommands {
		t.Errorf(
			"failed to assert that there were %d sub commands connected (was %d)",
			nbrOfCommands,
			len(commands),
		)
	}

	flags := cmd.PersistentFlags()
	flag := flags.Lookup(OldAccessTokenFlag)
	assert.NotNil(t, flag)
	assert.Equal(t, "t", flag.Shorthand)

	match := false
	viperKeys := viper.AllKeys()
	for _, key := range viperKeys {
		if key == AccessTokenFlag {
			match = true

			break
		}
	}
	assert.Truef(t, match, "failed to assert that flag was present: "+OldAccessTokenFlag)
	assert.Len(t, viperKeys, 23)
}

func TestPreRun(t *testing.T) {
	cmd := NewRootCmd("", wire.GetCliContainer())
	cmd.PreRun(cmd, nil)
}

func TestAccessTokenFlagDoesNotLeakEnvValueAsDefault(t *testing.T) {
	t.Setenv("DEBRICKED_TOKEN", "supersecrettoken")

	cmd := NewRootCmd("v0.0.0", wire.GetCliContainer())

	flag := cmd.PersistentFlags().Lookup(OldAccessTokenFlag)
	assert.NotNil(t, flag)
	assert.Empty(t, flag.DefValue, "flag default must not expose the env var value")
	assert.NotContains(t, cmd.UsageString(), "supersecrettoken")
}

func TestAccessTokenFallsBackToEnvValue(t *testing.T) {
	t.Setenv("DEBRICKED_TOKEN", "env-token")

	NewRootCmd("v0.0.0", wire.GetCliContainer())

	assert.Equal(t, "env-token", accessToken, "env var must still be used when the flag is omitted")
}

func TestAccessTokenFlagOverridesEnvValue(t *testing.T) {
	cases := []string{"--" + OldAccessTokenFlag, "-t"}
	for _, flagName := range cases {
		t.Run(flagName, func(t *testing.T) {
			t.Setenv("DEBRICKED_TOKEN", "env-token")

			cmd := NewRootCmd("v0.0.0", wire.GetCliContainer())
			output := &bytes.Buffer{}
			cmd.SetOut(output)
			cmd.SetErr(output)
			cmd.SetArgs([]string{flagName, "flag-token"})

			assert.NoError(t, cmd.Execute())
			assert.Equal(t, "flag-token", accessToken, "flag must take precedence over the env var")
			assert.NotContains(t, output.String(), "env-token", "help output must not contain the env var value")
		})
	}
}
