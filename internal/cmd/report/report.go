package report

import (
	"github.com/debricked/cli/internal/cmd/report/license"
	"github.com/debricked/cli/internal/cmd/report/sbom"
	"github.com/debricked/cli/internal/cmd/report/vulnerability"
	licenseReport "github.com/debricked/cli/internal/report/license"
	sbomReport "github.com/debricked/cli/internal/report/sbom"
	vulnerabilityReport "github.com/debricked/cli/internal/report/vulnerability"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewReportCmd(
	licenseReporter licenseReport.Reporter,
	vulnerabilityReporter vulnerabilityReport.Reporter,
	sbomReporter sbomReport.Reporter,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate reports",
		Long: `Generate reports.
Premium is required for license and vulnerability reports. Enterprise is required for SBOM reports. Please visit https://debricked.com/pricing/ for more info.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
	}

	cmd.AddCommand(license.NewLicenseCmd(licenseReporter))
	cmd.AddCommand(vulnerability.NewVulnerabilityCmd(vulnerabilityReporter))
	cmd.AddCommand(sbom.NewSBOMCmd(sbomReporter))

	return cmd
}
