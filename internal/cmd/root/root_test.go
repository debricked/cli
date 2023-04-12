package root

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd("v0.0.0")
	commands := cmd.Commands()
	nbrOfCommands := 3
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
	assert.Len(t, viperKeys, 13)
}

func TestPreRun(t *testing.T) {
	cmd := NewRootCmd("")
	cmd.PreRun(cmd, nil)
}
