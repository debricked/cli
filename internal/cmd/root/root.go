package root

import (
	"github.com/debricked/cli/internal/cmd/auth"
	"github.com/debricked/cli/internal/cmd/callgraph"
	"github.com/debricked/cli/internal/cmd/files"
	"github.com/debricked/cli/internal/cmd/fingerprint"
	"github.com/debricked/cli/internal/cmd/report"
	"github.com/debricked/cli/internal/cmd/resolve"
	"github.com/debricked/cli/internal/cmd/scan"
	"github.com/debricked/cli/internal/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var accessToken string

const AccessTokenFlag = "token"
const OldAccessTokenFlag = "access-token"

func NewRootCmd(version string, container *wire.CliContainer) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "debricked",
		Short: "Debricked CLI - Keep track of your dependencies!",
		Long: `A fast and flexible software composition analysis CLI tool, given to you by Debricked.
Complete documentation is available at https://docs.debricked.com/tools-and-integrations/cli/debricked-cli`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.PersistentFlags())
		},
		Version: version,
	}
	viper.Set("cliVersion", version)
	viper.SetEnvPrefix("DEBRICKED")
	viper.AutomaticEnv()
	viper.MustBindEnv(AccessTokenFlag)

	rootCmd.PersistentFlags().StringVarP(
		&accessToken,
		OldAccessTokenFlag,
		"t",
		viper.GetString(AccessTokenFlag),
		`Debricked access token. 
Read more: https://docs.debricked.com/product/administration/generate-access-token`,
	)

	var debClient = container.DebClient()
	debClient.SetAccessToken(&accessToken)

	rootCmd.AddCommand(report.NewReportCmd(container.LicenseReporter(), container.VulnerabilityReporter(), container.SBOMReporter()))
	rootCmd.AddCommand(files.NewFilesCmd(container.Finder()))
	rootCmd.AddCommand(scan.NewScanCmd(container.Scanner()))
	rootCmd.AddCommand(fingerprint.NewFingerprintCmd(container.Fingerprinter()))
	rootCmd.AddCommand(resolve.NewResolveCmd(container.Resolver()))
	rootCmd.AddCommand(callgraph.NewCallgraphCmd(container.CallgraphGenerator()))
	rootCmd.AddCommand(auth.NewAuthCmd(container.Authenticator()))

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	//rootCmd.SetVersionTemplate()

	return rootCmd
}
