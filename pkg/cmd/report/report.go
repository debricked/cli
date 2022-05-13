package report

import (
	"debricked/pkg/client"
	"debricked/pkg/cmd/report/license"
	"debricked/pkg/cmd/report/vulnerability"
	"github.com/spf13/cobra"
)

func NewReportCmd(debClient *client.DebClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate reports",
		Long: `Generate reports.
This is a premium feature. Please visit https://debricked.com/pricing/ for more info.`,
	}

	cmd.AddCommand(license.NewLicenseCmd(debClient))
	cmd.AddCommand(vulnerability.NewVulnerabilityCmd(debClient))

	return cmd
}
