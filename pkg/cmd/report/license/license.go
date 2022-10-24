package license

import (
	"errors"
	"fmt"
	"github.com/debricked/cli/pkg/report"
	"github.com/debricked/cli/pkg/report/license"
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
		RunE: RunE(reporter),
	}

	cmd.Flags().StringVarP(&email, EmailFlag, "e", "", "The email address that the report will be sent to")
	_ = cmd.MarkFlagRequired(EmailFlag)
	viper.MustBindEnv(EmailFlag)

	cmd.Flags().StringVarP(&commitHash, CommitFlag, "c", "", "commit hash")
	_ = cmd.MarkFlagRequired(CommitFlag)
	viper.MustBindEnv(CommitFlag)

	_ = viper.BindPFlags(cmd.Flags())

	return cmd
}

func RunE(r report.IReporter) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		orderArgs := license.OrderArgs{
			Email:      viper.GetString(EmailFlag),
			CommitHash: viper.GetString(CommitFlag),
		}

		if err := r.Order(orderArgs); err != nil {
			return errors.New(fmt.Sprintf("%s %s\n", color.RedString("⨯"), err.Error()))
		}

		fmt.Printf("%s Successfully ordered license report\n", color.GreenString("✔"))

		return nil
	}
}
