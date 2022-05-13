package report

import (
	"debricked/pkg/client"
	"debricked/pkg/cmd/report/license"
	"github.com/spf13/cobra"
)

var debClient *client.DebClient

func NewReportCmd(debrickedClient *client.DebClient) *cobra.Command {
	debClient = debrickedClient
	cmd := &cobra.Command{
		Use:   "report",
		Short: "Generate reports",
		Long: `Generate reports for a commit.
This is a premium feature. Please visit https://debricked.com/pricing/ for more info.`,
	}

	cmd.AddCommand(license.NewLicenseCmd(debClient))

	return cmd
}
