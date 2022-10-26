package report

import (
	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/cmd/report/license"
	"github.com/debricked/cli/pkg/cmd/report/vulnerability"
	licenseReport "github.com/debricked/cli/pkg/report/license"
	vulnerabilityReport "github.com/debricked/cli/pkg/report/vulnerability"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewReportCmd(debClient *client.IDebClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate reports",
		Long: `Generate reports.
This is a premium feature. Please visit https://debricked.com/pricing/ for more info.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
	}

	lReporter := licenseReport.Reporter{DebClient: *debClient}
	cmd.AddCommand(license.NewLicenseCmd(lReporter))

	vReporter := vulnerabilityReport.Reporter{DebClient: *debClient}
	cmd.AddCommand(vulnerability.NewVulnerabilityCmd(vReporter))

	return cmd
}
