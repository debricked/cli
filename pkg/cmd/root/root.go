package root

import (
	"debricked/pkg/client"
	"debricked/pkg/cmd/check"
	"debricked/pkg/cmd/files"
	"debricked/pkg/cmd/login"
	"debricked/pkg/cmd/report"
	"debricked/pkg/cmd/scan"
	"github.com/spf13/cobra"
)

var accessToken string

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "debricked",
		Short: "Debricked CLI - Keep track of your dependencies!",
		Long: `A fast and flexible software composition analysis CLI tool, given to you by Debricked.
Complete documentation is available at https://debricked.com/docs/integrations/cli.html#debricked-cli`,
	}

	rootCmd.PersistentFlags().StringVarP(
		&accessToken,
		"access-token",
		"t",
		"",
		`Debricked access token. 
Read more: https://debricked.com/docs/administration/access-tokens.html`,
	)

	var debClient client.Client = client.NewDebClient(&accessToken)

	rootCmd.AddCommand(report.NewReportCmd(&debClient))
	rootCmd.AddCommand(scan.NewScanCmd(&debClient))
	rootCmd.AddCommand(check.NewCheckCmd(&debClient))
	rootCmd.AddCommand(login.NewLoginCmd(&debClient))
	rootCmd.AddCommand(files.NewFilesCmd(&debClient))

	return rootCmd
}
