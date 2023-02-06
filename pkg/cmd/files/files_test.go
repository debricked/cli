package files

import (
	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/client/testdata"
	"testing"
)

func TestNewFilesCmd(t *testing.T) {
	var debClient client.IDebClient = testdata.NewDebClientMock()
	cmd := NewFilesCmd(&debClient)
	commands := cmd.Commands()
	nbrOfCommands := 1
	if len(commands) != nbrOfCommands {
		t.Errorf("failed to assert that there were %d sub commands connected", nbrOfCommands)
	}
}
