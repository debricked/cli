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
	nbrOfCommands := 5
	if len(commands) != nbrOfCommands {
		t.Errorf("failed to assert that there were %d sub commands connected", nbrOfCommands)
	}

	flags := cmd.PersistentFlags()
	flag := flags.Lookup(AccessTokenFlag)
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
	assert.Truef(t, match, "failed to assert that flag was present: "+AccessTokenFlag)
	assert.Len(t, viperKeys, 14)
}

func TestPreRun(t *testing.T) {
	cmd := NewRootCmd("", wire.GetCliContainer())
	cmd.PreRun(cmd, nil)
}
