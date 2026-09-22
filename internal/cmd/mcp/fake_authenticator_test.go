package mcp

import (
	"golang.org/x/oauth2"
)

// fakeAuthenticator is a test double for auth.IAuthenticator.
type fakeAuthenticator struct {
	token *oauth2.Token
	err   error
}

func (f fakeAuthenticator) Authenticate() error { return nil }
func (f fakeAuthenticator) Logout() error       { return nil }
func (f fakeAuthenticator) Token() (*oauth2.Token, error) {
	return f.token, f.err
}
