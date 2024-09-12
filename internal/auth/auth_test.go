package auth

import (
	"testing"

	"github.com/debricked/cli/internal/auth/testdata"
	clientTestdata "github.com/debricked/cli/internal/client/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
)

func TestNewAuthenticator(t *testing.T) {
	res := NewDebrickedAuthenticator(clientTestdata.NewDebClientMock())
	assert.NotNil(t, res)
}

func TestSecretClientSet(t *testing.T) {
	user := "TestDebrickedCLIUserSet"
	service := "TestDebrickedCLIServiceSet"
	secret := "TestDebrickedCLISecretSet"
	dsc := DebrickedSecretClient{user}
	_, err := keyring.Get(service, user)
	assert.Error(t, err)
	err = dsc.Set(service, secret)
	assert.NoError(t, err)
	savedSecret, err := keyring.Get(service, user)
	assert.NoError(t, err)
	err = keyring.Delete(service, user)
	assert.NoError(t, err)
	assert.Equal(t, secret, savedSecret)
}

func TestSecretClientGet(t *testing.T) {
	user := "TestDebrickedCLIUserGet"
	service := "TestDebrickedCLIServiceGet"
	secret := "TestDebrickedCLISecretGet"
	dsc := DebrickedSecretClient{user}
	err := keyring.Set(service, user, secret)
	assert.NoError(t, err)
	savedSecret, err := dsc.Get(service)
	assert.NoError(t, err)
	err = keyring.Delete(service, user)
	assert.NoError(t, err)
	assert.Equal(t, secret, savedSecret)
}

func TestSecretClientDelete(t *testing.T) {
	user := "TestDebrickedCLIUserDelete"
	service := "TestDebrickedCLIServiceDelete"
	secret := "TestDebrickedCLISecretDelete"
	dsc := DebrickedSecretClient{user}
	err := keyring.Set(service, user, secret)
	assert.NoError(t, err)
	savedSecret, err := keyring.Get(service, user)
	assert.NoError(t, err)
	assert.Equal(t, secret, savedSecret)
	err = dsc.Delete(service)
	assert.NoError(t, err)
	_, err = keyring.Get(service, user)
	assert.Error(t, err)
}

func TestMockedLogout(t *testing.T) {
	authenticator := Authenticator{
		SecretClient: testdata.MockSecretClient{},
		OAuthConfig:  nil,
	}
	err := authenticator.Logout()

	assert.NoError(t, err)
}

func TestMockedLogoutErrorRefresh(t *testing.T) {
	authenticator := Authenticator{
		SecretClient: testdata.MockErrorSecretClient{
			ErrorPattern: "Refresh",
		},
		OAuthConfig: nil,
	}
	err := authenticator.Logout()

	assert.Error(t, err)
}

func TestMockedSaveToken(t *testing.T) {
	authenticator := Authenticator{
		SecretClient: testdata.MockSecretClient{},
		OAuthConfig:  nil,
	}
	token := &oauth2.Token{
		RefreshToken: "refreshToken",
		AccessToken:  "accessToken",
	}
	err := authenticator.saveToken(token)

	assert.NoError(t, err)
}

func TestMockedSaveTokenRefreshError(t *testing.T) {
	authenticator := Authenticator{
		SecretClient: testdata.MockErrorSecretClient{
			ErrorPattern: "Refresh",
		},
		OAuthConfig: nil,
	}
	token := &oauth2.Token{
		RefreshToken: "refreshToken",
		AccessToken:  "accessToken",
	}
	err := authenticator.saveToken(token)

	assert.Error(t, err)
}

func TestMockedToken(t *testing.T) {
	authenticator := Authenticator{
		SecretClient: testdata.MockSecretClient{},
		OAuthConfig:  nil,
	}
	token, err := authenticator.Token()

	assert.NoError(t, err)
	assert.Equal(t, token.TokenType, "jwt")
	assert.Equal(t, token.RefreshToken, "token")
	assert.Equal(t, token.AccessToken, "token")
}

func TestMockedTokenRefreshError(t *testing.T) {
	authenticator := Authenticator{
		SecretClient: testdata.MockErrorSecretClient{
			ErrorPattern: "Refresh",
		},
		OAuthConfig: nil,
	}
	_, err := authenticator.Token()

	assert.Error(t, err)
}

func TestMockedTokenAccessError(t *testing.T) {
	authenticator := Authenticator{
		SecretClient: testdata.MockErrorSecretClient{
			ErrorPattern: "Access",
		},
		OAuthConfig: nil,
	}
	_, err := authenticator.Token()

	assert.Error(t, err)
}

func TestMockedAuthenticate(t *testing.T) {
	authenticator := Authenticator{
		SecretClient:  testdata.MockSecretClient{},
		OAuthConfig:   testdata.MockOAuthConfig{},
		AuthWebHelper: testdata.MockAuthWebHelper{},
	}
	err := authenticator.Authenticate()

	assert.NoError(t, err)
}

func TestMockedAuthenticateOpenURLError(t *testing.T) {
	authenticator := Authenticator{
		SecretClient:  testdata.MockSecretClient{},
		OAuthConfig:   testdata.MockOAuthConfig{},
		AuthWebHelper: testdata.MockErrorAuthWebHelper{},
	}
	err := authenticator.Authenticate()

	assert.Error(t, err)
}
