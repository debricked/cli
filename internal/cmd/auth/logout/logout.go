package logout

import (
	"fmt"
	"github.com/debricked/cli/internal/auth"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewLogoutCmd(authenticator auth.IAuthenticator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Logout debricked user",
		Long:  `Remove cached credentials to logout debricked user.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(authenticator),
	}

	return cmd
}

func RunE(a auth.IAuthenticator) func(_ *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		err := a.Logout()
		if err != nil {
			return err
		}
		fmt.Printf(
			"%s Successfully removed credentials",
			color.GreenString("✔"),
		)

		return nil
	}
}
