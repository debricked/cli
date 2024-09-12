package testdata

import (
	"context"
	"strings"

	"golang.org/x/oauth2"
)

type MockError struct{}

func (me MockError) Error() string {
	return "MockError!"
}

type MockSecretClient struct{}

func (msc MockSecretClient) Set(service, secret string) error {
	return nil
}

func (msc MockSecretClient) Get(service string) (string, error) {
	return "token", nil
}

func (msc MockSecretClient) Delete(service string) error {
	return nil
}

type MockErrorSecretClient struct {
	ErrorPattern string
}

func (msc MockErrorSecretClient) Set(service, secret string) error {
	if strings.Contains(service, msc.ErrorPattern) {

		return MockError{}
	}

	return nil
}

func (msc MockErrorSecretClient) Get(service string) (string, error) {
	if strings.Contains(service, msc.ErrorPattern) {

		return "", MockError{}
	}
	return "token", nil
}

func (msc MockErrorSecretClient) Delete(service string) error {
	if strings.Contains(service, msc.ErrorPattern) {

		return MockError{}
	}
	return nil
}

type MockAuthenticator struct{}

type ErrorMockAuthenticator struct{}

type MockOAuthConfig struct{}

type MockAuthWebHelper struct{}

type MockErrorAuthWebHelper struct{}

func (ma MockAuthenticator) Authenticate() error {
	return nil
}

func (ma MockAuthenticator) Logout() error {
	return nil
}

func (ma MockAuthenticator) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		RefreshToken: "refresh",
		AccessToken:  "access",
		TokenType:    "jwt",
	}, nil
}

func (ma ErrorMockAuthenticator) Authenticate() error {
	return MockError{}
}

func (ma ErrorMockAuthenticator) Logout() error {
	return MockError{}
}

func (ma ErrorMockAuthenticator) Token() (*oauth2.Token, error) {
	return nil, MockError{}
}

func (mawh MockAuthWebHelper) OpenURL(string) error {
	return nil
}

func (mawh MockAuthWebHelper) Callback(string) string {
	return "callback"
}

func (mawh MockErrorAuthWebHelper) OpenURL(string) error {
	return MockError{}
}

func (mawh MockErrorAuthWebHelper) Callback(string) string {
	return "callback"
}

func (moc MockOAuthConfig) Exchange(context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken:  "accessToken",
		RefreshToken: "accessToken",
	}, nil
}

func (moc MockOAuthConfig) AuthCodeURL(string, ...oauth2.AuthCodeOption) string {
	return "localhost"
}
