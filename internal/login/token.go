package login

import (
	"github.com/zalando/go-keyring"

	"golang.org/x/oauth2"
)

type SecretClient interface {
	Set(string, string) error
	Get(string) (string, error)
}

type DebrickedSecretClient struct {
	User string
}

type DebrickedTokenSource struct {
	SecretClient SecretClient
}

func (dsc DebrickedSecretClient) Set(service, secret string) error {
	return keyring.Set(service, dsc.User, secret)
}

func (dsc DebrickedSecretClient) Get(service string) (string, error) {
	return keyring.Get(service, dsc.User)
}

func GetDebrickedTokenSource() oauth2.TokenSource {
	return DebrickedTokenSource{
		SecretClient: DebrickedSecretClient{
			User: "DebrickedCLI",
		},
	}
}

func (dts DebrickedTokenSource) Token() (*oauth2.Token, error) {
	refreshToken, err := dts.SecretClient.Get("DebrickedRefreshToken")
	if err != nil {
		if err == keyring.ErrNotFound {
			// refreshToken is not yet set, initialize authorization
			authenticator := Authenticator{
				ClientID: "01919462-7d6e-78e8-aa24-ba779213c90f",
				Scopes:   []string{"select", "profile", "basicRepo"},
			}
			token, err := authenticator.Authenticate()
			if err != nil {
				return nil, err
			}
			dts.SecretClient.Set("DebrickedRefreshToken", token.RefreshToken)
			dts.SecretClient.Set("DebrickedAccessToken", token.AccessToken)
		} else {
			return nil, err
		}
	}
	accessToken, err := dts.SecretClient.Get("DebrickedAccessToken")
	if err != nil {
		accessToken = ""
	}
	// TODO: Parse expiry date
	return &oauth2.Token{
		RefreshToken: refreshToken,
		TokenType:    "jwt",
		AccessToken:  accessToken,
	}, nil
}
