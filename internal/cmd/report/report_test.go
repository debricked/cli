package report

import (
	"testing"

	"github.com/debricked/cli/internal/report/license"
	"github.com/debricked/cli/internal/report/sbom"
	"github.com/debricked/cli/internal/report/vulnerability"
	"github.com/stretchr/testify/assert"
)

func TestNewReportCmd(t *testing.T) {
	cmd := NewReportCmd(license.Reporter{}, vulnerability.Reporter{}, sbom.Reporter{})
	commands := cmd.Commands()
	nbrOfCommands := 3
	assert.Lenf(t, commands, nbrOfCommands, "failed to assert that there were %d sub commands connected", nbrOfCommands)
}

func TestPreRun(t *testing.T) {
	var licenseReporter license.Reporter
	var vulnReporter vulnerability.Reporter
	var sbomReporter sbom.Reporter
	cmd := NewReportCmd(licenseReporter, vulnReporter, sbomReporter)
	cmd.PreRun(cmd, nil)
}
