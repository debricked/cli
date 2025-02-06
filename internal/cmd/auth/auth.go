package auth

import (
	"github.com/debricked/cli/internal/auth"
	"github.com/debricked/cli/internal/cmd/auth/login"
	"github.com/debricked/cli/internal/cmd/auth/logout"
	"github.com/debricked/cli/internal/cmd/auth/token"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewAuthCmd(authenticator auth.IAuthenticator) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Debricked authentication.",
		Long:  `Debricked service authentication. Currently in beta and will most likely not work as expected`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		Hidden: true,
	}
	cmd.AddCommand(login.NewLoginCmd(authenticator))
	cmd.AddCommand(logout.NewLogoutCmd(authenticator))
	cmd.AddCommand(token.NewTokenCmd(authenticator))

	return cmd
}
