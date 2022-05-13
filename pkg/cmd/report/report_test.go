package report

import (
	"debricked/pkg/client"
	"fmt"
	"testing"
)

func TestNewReportCmd(t *testing.T) {
	cmd := NewReportCmd(client.NewDebClient(nil))
	commands := cmd.Commands()
	nbrOfCommands := 2
	if len(commands) != nbrOfCommands {
		t.Error(fmt.Sprintf("failed to assert that there were %d sub commands connected", nbrOfCommands))
	}
}
