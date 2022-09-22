package root

import (
	"fmt"
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
	flag := flags.Lookup("access-token")
	if flag == nil {
		t.Error("failed to assert that access-token flag was set")
	}
	if flag.Shorthand != "t" {
		t.Error("failed to assert that access-token flag shorthand was set correctly")
	}
}
