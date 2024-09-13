package testdata

import (
	"context"
	"strings"

	"golang.org/x/oauth2"
)

type MockError struct {
	Message string
}

func (me MockError) Error() string {
	return me.Message
}

type MockSecretClient struct{}

type MockExpiredSecretClient struct{}

type MockInvalidSecretClient struct{}

func (msc MockSecretClient) Set(service, secret string) error {
	return nil
}

func (msc MockSecretClient) Get(service string) (string, error) {
	return "token", nil
}

func (msc MockSecretClient) Delete(service string) error {
	return nil
}

func (msc MockExpiredSecretClient) Set(service, secret string) error {
	return nil
}

func (msc MockExpiredSecretClient) Get(service string) (string, error) {
	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiIwMTkxOTQ2Mi03ZDZlLTc4ZTgtYWEyNC1iYTc3OTIxM2M5MGYiLCJqdGkiOiJlMTdhMmFlYTk0ZjgyNTdjYWU1NWM3ZjRiNTczNTRiMzI2YmNiYTZiZmY3ZGQ0ZWQ2NjU3NDA4MWE4ODFjN2VhMmM3OGU3Y2EzM2UxMjU5MyIsImlhdCI6MTY5NDU5NzkzNy4zNjAwMTUsIm5iZiI6MTY5NDU5NzkzNy4zNjAwMTcsImV4cCI6MTY5NDU5NzkzNy4zNTM3MDMsInN1YiI6ImZpbGlwLmhlZGVuK2FkbWluQGRlYnJpY2tlZC5jb20iLCJzY29wZXMiOlsic2VsZWN0IiwicHJvZmlsZSIsImJhc2ljUmVwbyJdfQ.CMqnQM9QFHTthDMv4K8q6gmkkFmbOIhrmKXwfo7kMWU", nil
}

func (msc MockExpiredSecretClient) Delete(service string) error {
	return nil
}

func (msc MockInvalidSecretClient) Set(service, secret string) error {
	return nil
}

func (msc MockInvalidSecretClient) Get(service string) (string, error) {
	return "eyJhdWQiOiIwMTkxOTQ2Mi03ZDZlLTc4ZTgtYWEyNC1iYTc3OTIxM2M5MGYiLCJqdGkiOiJlMTdhMmFlYTk0ZjgyNTdjYWU1NWM3ZjRiNTczNTRiMzI2YmNiYTZiZmY3ZGQ0ZWQ2NjU3NDA4MWE4ODFjN2VhMmM3OGU3Y2EzM2UxMjU5MyIsImlhdCI6MTY5NDU5NzkzNy4zNjAwMTUsIm5iZiI6MTY5NDU5NzkzNy4zNjAwMTcsImV4cCI6MTY5NDU5NzkzNy4zNTM3MDMsInN1YiI6ImZpbGlwLmhlZGVuK2FkbWluQGRlYnJpY2tlZC5jb20iLCJzY29wZXMiOlsic2VsZWN0IiwicHJvZmlsZSIsImJhc2ljUmVwbyJdfQ", nil
}

func (msc MockInvalidSecretClient) Delete(service string) error {
	return nil
}

type MockErrorSecretClient struct {
	ErrorPattern string
	Message      string
}

func (msc MockErrorSecretClient) Set(service, secret string) error {
	if strings.Contains(service, msc.ErrorPattern) {

		return MockError{
			Message: msc.Message,
		}
	}

	return nil
}

func (msc MockErrorSecretClient) Get(service string) (string, error) {
	if strings.Contains(service, msc.ErrorPattern) {

		return "", MockError{
			Message: msc.Message,
		}
	}
	return "token", nil
}

func (msc MockErrorSecretClient) Delete(service string) error {
	if strings.Contains(service, msc.ErrorPattern) {

		return MockError{
			Message: msc.Message,
		}
	}
	return nil
}

type MockAuthenticator struct{}

type ErrorMockAuthenticator struct{}

type MockOAuthConfig struct {
	MockTokenSource oauth2.TokenSource
}

type MockOAuthConfigExchangeError struct{}

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
	return MockError{""}
}

func (ma ErrorMockAuthenticator) Logout() error {
	return MockError{""}
}

func (ma ErrorMockAuthenticator) Token() (*oauth2.Token, error) {
	return nil, MockError{""}
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

func (moc MockOAuthConfigExchangeError) Exchange(context.Context, string, ...oauth2.AuthCodeOption) (*oauth2.Token, error) {
	return nil, MockError{Message: "HTTP Error"}
}

func (moc MockOAuthConfigExchangeError) AuthCodeURL(string, ...oauth2.AuthCodeOption) string {
	return "localhost"
}

func (moc MockOAuthConfigExchangeError) TokenSource(context.Context, *oauth2.Token) oauth2.TokenSource {
	return nil
}

type MockTokenSource struct {
	StaticToken *oauth2.Token
	Error       error
}

func (mts MockTokenSource) Token() (*oauth2.Token, error) {
	if mts.Error != nil {
		return nil, mts.Error
	}
	return mts.StaticToken, nil
}

func (moc MockOAuthConfig) TokenSource(context.Context, *oauth2.Token) oauth2.TokenSource {
	return moc.MockTokenSource
}
