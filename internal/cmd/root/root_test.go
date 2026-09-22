package root

import (
	"testing"

	"github.com/debricked/cli/internal/wire"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd("v0.0.0", wire.GetCliContainer())
	commands := cmd.Commands()
	nbrOfCommands := 9
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

	mcpMatch := false
	for _, subCmd := range commands {
		if subCmd.Name() == "mcp" {
			mcpMatch = true

			break
		}
	}
	assert.True(t, mcpMatch, "failed to assert that mcp command was registered")
}

func TestPreRun(t *testing.T) {
	cmd := NewRootCmd("", wire.GetCliContainer())
	cmd.PreRun(cmd, nil)
}
