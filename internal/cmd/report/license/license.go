package license

import (
	"fmt"

	"github.com/debricked/cli/internal/report"
	"github.com/debricked/cli/internal/report/license"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var email string
var commitHash string

const (
	EmailFlag  = "email"
	CommitFlag = "commit"
)

func NewLicenseCmd(reporter report.IReporter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "license",
		Short: "Generate license report",
		Long: `Generate license report from a commit hash. 
This is a premium feature. Please visit https://debricked.com/pricing/ for more info.
The finished report will be sent to the specified email address.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(reporter),
	}

	cmd.Flags().StringVarP(&email, EmailFlag, "e", "", "The email address that the report will be sent to")
	viper.MustBindEnv(EmailFlag)

	cmd.Flags().StringVarP(&commitHash, CommitFlag, "c", "", "commit hash")
	viper.MustBindEnv(CommitFlag)

	return cmd
}

func RunE(r report.IReporter) func(_ *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		orderArgs := license.OrderArgs{
			Email:      viper.GetString(EmailFlag),
			CommitHash: viper.GetString(CommitFlag),
		}

		if err := r.Order(orderArgs); err != nil {
			return fmt.Errorf("%s %s\n", color.RedString("⨯"), err.Error())
		}

		fmt.Printf("%s Successfully ordered license report\n", color.GreenString("✔"))

		return nil
	}
}
