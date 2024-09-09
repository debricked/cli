package token

import (
	"fmt"
	"github.com/debricked/cli/internal/auth"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewTokenCmd(authenticator auth.IAuthenticator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Retrieve access token",
		Long:  `Retrieve access token for currently logged in Debricked user.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(authenticator),
	}

	return cmd
}

func RunE(a auth.IAuthenticator) func(_ *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		token, err := a.Token()
		if err != nil {
			return err
		}
		fmt.Printf(
			"Refresh Token = %s\nAccess Token = %s\n",
			color.BlueString(token.RefreshToken),
			color.BlueString(token.AccessToken),
		)

		return nil
	}
}
