package report

import (
	"github.com/debricked/cli/pkg/client"
	"testing"
)

func TestNewReportCmd(t *testing.T) {
	var c client.IDebClient = client.NewDebClient(nil)
	cmd := NewReportCmd(&c)
	commands := cmd.Commands()
	nbrOfCommands := 2
	if len(commands) != nbrOfCommands {
		t.Errorf("failed to assert that there were %d sub commands connected", nbrOfCommands)
	}
}
