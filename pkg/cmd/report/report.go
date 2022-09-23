package report

import (
	"debricked/pkg/client"
	"debricked/pkg/cmd/report/license"
	"debricked/pkg/cmd/report/vulnerability"
	licenseReport "debricked/pkg/report/license"
	vulnerabilityReport "debricked/pkg/report/vulnerability"
	"github.com/spf13/cobra"
)

func NewReportCmd(debClient *client.IDebClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate reports",
		Long: `Generate reports.
This is a premium feature. Please visit https://debricked.com/pricing/ for more info.`,
	}

	lReporter := licenseReport.Reporter{DebClient: *debClient}
	cmd.AddCommand(license.NewLicenseCmd(lReporter))

	vReporter := vulnerabilityReport.Reporter{DebClient: *debClient}
	cmd.AddCommand(vulnerability.NewVulnerabilityCmd(vReporter))

	return cmd
}
