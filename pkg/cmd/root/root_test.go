package root

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"
)

func TestNewRootCmd(t *testing.T) {
	cmd := NewRootCmd()
	commands := cmd.Commands()
	nbrOfCommands := 3
	if len(commands) != nbrOfCommands {
		t.Error(fmt.Sprintf("failed to assert that there were %d sub commands connected", nbrOfCommands))
	}

	flags := cmd.PersistentFlags()
	flag := flags.Lookup(AccessTokenFlag)
	if flag == nil {
		t.Error("failed to assert that access-token flag was set")
	}
	if flag.Shorthand != "t" {
		t.Error("failed to assert that access-token flag shorthand was set correctly")
	}

	match := false
	viperKeys := viper.AllKeys()
	for _, key := range viperKeys {
		if key == AccessTokenFlag {
			match = true
			break
		}
	}
	if !match {
		t.Error("failed to assert that flag was present: " + AccessTokenFlag)
	}

	if len(viperKeys) != 11 {
		t.Error("failed to assert number of keys bound by Viper")
	}

}
