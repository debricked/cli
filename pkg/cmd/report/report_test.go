package report

import (
	"testing"

	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/client/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewReportCmd(t *testing.T) {
	var c client.IDebClient = testdata.NewDebClientMock()
	cmd := NewReportCmd(&c)
	commands := cmd.Commands()
	nbrOfCommands := 2
	assert.Lenf(t, commands, nbrOfCommands, "failed to assert that there were %d sub commands connected", nbrOfCommands)
}

func TestPreRun(t *testing.T) {
	var c client.IDebClient = testdata.NewDebClientMock()
	cmd := NewReportCmd(&c)
	cmd.PreRun(cmd, nil)
}
