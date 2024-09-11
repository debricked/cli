package token

import (
	"encoding/json"
	"fmt"

	"github.com/debricked/cli/internal/auth"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var jsonFormat bool

const JsonFlag = "json"

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
	cmd.Flags().BoolVarP(&jsonFormat, JsonFlag, "j", false, `Print files in JSON format
Format:
[
  {
    "access_token": <access token>,
    "token_type": "jwt",
    "refresh_token": <refresh token>,
    "expiry": <access token expiry date>,
  },
]
`)
	viper.MustBindEnv(JsonFlag)

	return cmd
}

func RunE(a auth.IAuthenticator) func(_ *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		token, err := a.Token()
		if err != nil {
			return err
		}
		if viper.GetBool(JsonFlag) {
			jsonToken, _ := json.Marshal(token)
			fmt.Println(string(jsonToken))
		} else {
			fmt.Printf(
				"Refresh Token = %s\nAccess Token = %s\n",
				color.BlueString(token.RefreshToken),
				color.BlueString(token.AccessToken),
			)
		}

		return nil
	}
}
