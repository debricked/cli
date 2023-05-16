package report

import (
	"github.com/debricked/cli/internal/cmd/report/license"
	"github.com/debricked/cli/internal/cmd/report/vulnerability"
	licenseReport "github.com/debricked/cli/internal/report/license"
	vulnerabilityReport "github.com/debricked/cli/internal/report/vulnerability"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewReportCmd(
	licenseReporter licenseReport.Reporter,
	vulnerabilityReporter vulnerabilityReport.Reporter,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate reports",
		Long: `Generate reports.
This is a premium feature. Please visit https://debricked.com/pricing/ for more info.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
	}

	cmd.AddCommand(license.NewLicenseCmd(licenseReporter))
	cmd.AddCommand(vulnerability.NewVulnerabilityCmd(vulnerabilityReporter))

	return cmd
}
