package login

import (
	"fmt"
	"github.com/debricked/cli/internal/login"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewLoginCmd(authenticator login.IAuthenticator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate debricked user",
		Long:  `Start authentication flow to generate access token.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(authenticator),
	}

	return cmd
}

func RunE(a login.IAuthenticator) func(_ *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		token, err := a.Authenticate()
		if err != nil {
			return err
		}
		fmt.Printf(
			"%s Successfully authenticated\nToken=%s",
			color.GreenString("✔"),
			color.BlueString(token),
		)

		return nil
	}
}
