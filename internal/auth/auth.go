package auth

import (
	"context"
	"github.com/debricked/cli/internal/client"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

type IAuthenticator interface {
	Authenticate() error
	Logout() error
	Token() (*oauth2.Token, error)
}

type ISecretClient interface {
	Set(string, string) error
	Get(string) (string, error)
	Delete(string) error
}

type IOAuthConfig interface {
	AuthCodeURL(string, ...oauth2.AuthCodeOption) string
	Exchange(context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error)
} // Wrapping interface for config to enable mocking

type Authenticator struct {
	SecretClient  ISecretClient
	OAuthConfig   IOAuthConfig
	AuthWebHelper IAuthWebHelper
}

type DebrickedSecretClient struct {
	User string
}

func (dsc DebrickedSecretClient) Set(service, secret string) error {
	return keyring.Set(service, dsc.User, secret)
}

func (dsc DebrickedSecretClient) Get(service string) (string, error) {
	return keyring.Get(service, dsc.User)
}

func (dsc DebrickedSecretClient) Delete(service string) error {
	return keyring.Delete(service, dsc.User)
}

func NewDebrickedAuthenticator(client client.IDebClient) Authenticator {
	return Authenticator{
		SecretClient: DebrickedSecretClient{
			User: "DebrickedCLI",
		},
		OAuthConfig: &oauth2.Config{
			ClientID:     "01919462-7d6e-78e8-aa24-ba779213c90f",
			ClientSecret: "",
			Endpoint: oauth2.Endpoint{
				AuthURL:  client.Host() + "/app/oauth/authorize",
				TokenURL: client.Host() + "/app/oauth/token",
			},
			RedirectURL: "http://localhost:9096/callback",
			Scopes:      []string{"select", "profile", "basicRepo"},
		},
		AuthWebHelper: NewAuthWebHelper(),
	}
}

func (a Authenticator) Logout() error {
	err := a.SecretClient.Delete("DebrickedRefreshToken")
	if err != nil {
		return err
	}
	err = a.SecretClient.Delete("DebrickedAccessToken")
	return err
}

func (a Authenticator) Token() (*oauth2.Token, error) {
	refreshToken, err := a.SecretClient.Get("DebrickedRefreshToken")
	if err != nil {
		return nil, err
	}
	accessToken, err := a.SecretClient.Get("DebrickedAccessToken")
	if err != nil {
		return nil, err
	}

	return &oauth2.Token{
		RefreshToken: refreshToken,
		TokenType:    "jwt",
		AccessToken:  accessToken,
	}, nil
}

func (a Authenticator) saveToken(token *oauth2.Token) error {
	err := a.SecretClient.Set("DebrickedRefreshToken", token.RefreshToken)
	if err != nil {
		return err
	}
	err = a.SecretClient.Set("DebrickedAccessToken", token.AccessToken)

	return err
}

func (a Authenticator) exchange(authCode, codeVerifier string) (*oauth2.Token, error) {

	return a.OAuthConfig.Exchange(
		context.Background(),
		authCode,
		oauth2.VerifierOption(codeVerifier),
	)

}

func (a Authenticator) Authenticate() error {
	state := oauth2.GenerateVerifier()
	codeVerifier := oauth2.GenerateVerifier()
	authURL := a.OAuthConfig.AuthCodeURL(
		state,
		oauth2.S256ChallengeOption(codeVerifier),
	)

	err := a.AuthWebHelper.OpenURL(authURL)
	if err != nil {
		return err
	}

	authCode := a.AuthWebHelper.Callback(state)
	token, err := a.exchange(authCode, codeVerifier)
	if err != nil {
		return err
	}

	return a.saveToken(token)
}
