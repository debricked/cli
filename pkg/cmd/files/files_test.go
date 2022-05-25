package files

import (
	"debricked/pkg/client"
	"fmt"
	"testing"
)

func TestNewFilesCmd(t *testing.T) {
	cmd := NewFilesCmd(client.NewDebClient(nil))
	commands := cmd.Commands()
	nbrOfCommands := 1
	if len(commands) != nbrOfCommands {
		t.Error(fmt.Sprintf("failed to assert that there were %d sub commands connected", nbrOfCommands))
	}
}
