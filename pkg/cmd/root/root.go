package root

import (
	"github.com/debricked/cli/pkg/cmd/files"
	"github.com/debricked/cli/pkg/cmd/report"
	"github.com/debricked/cli/pkg/cmd/resolve"
	"github.com/debricked/cli/pkg/cmd/scan"
	"github.com/debricked/cli/pkg/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accessToken string

const AccessTokenFlag = "access-token"

func NewRootCmd(version string, container *wire.CliContainer) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "debricked",
		Short: "Debricked CLI - Keep track of your dependencies!",
		Long: `A fast and flexible software composition analysis CLI tool, given to you by Debricked.
Complete documentation is available at https://debricked.com/docs/integrations/cli.html#debricked-cli`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.PersistentFlags())
		},
		Version: version,
	}
	viper.SetEnvPrefix("DEBRICKED")
	viper.MustBindEnv(AccessTokenFlag)
	rootCmd.PersistentFlags().StringVarP(
		&accessToken,
		AccessTokenFlag,
		"t",
		"",
		`Debricked access token. 
Read more: https://debricked.com/docs/administration/access-tokens.html`,
	)

	var debClient = container.DebClient()
	debClient.SetAccessToken(&accessToken)

	rootCmd.AddCommand(report.NewReportCmd(container.LicenseReporter(), container.VulnerabilityReporter()))
	rootCmd.AddCommand(files.NewFilesCmd(container.Finder()))
	rootCmd.AddCommand(scan.NewScanCmd(container.Scanner()))
	rootCmd.AddCommand(resolve.NewResolveCmd(container.Resolver()))

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	//rootCmd.SetVersionTemplate()

	return rootCmd
}
