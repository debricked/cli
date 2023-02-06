package report

import (
	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/client/testdata"
	"testing"
)

func TestNewReportCmd(t *testing.T) {
	var c client.IDebClient = testdata.NewDebClientMock()
	cmd := NewReportCmd(&c)
	commands := cmd.Commands()
	nbrOfCommands := 2
	if len(commands) != nbrOfCommands {
		t.Errorf("failed to assert that there were %d sub commands connected", nbrOfCommands)
	}
}
