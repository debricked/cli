package report

import (
	"testing"

	"github.com/debricked/cli/pkg/report/license"
	"github.com/debricked/cli/pkg/report/vulnerability"
	"github.com/stretchr/testify/assert"
)

func TestNewReportCmd(t *testing.T) {
	cmd := NewReportCmd(license.Reporter{}, vulnerability.Reporter{})
	commands := cmd.Commands()
	nbrOfCommands := 2
	assert.Lenf(t, commands, nbrOfCommands, "failed to assert that there were %d sub commands connected", nbrOfCommands)
}

func TestPreRun(t *testing.T) {
	var c client.IDebClient = testdata.NewDebClientMock()
	cmd := NewReportCmd(&c)
	cmd.PreRun(cmd, nil)
}
