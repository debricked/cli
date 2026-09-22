package mcp

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/debricked/cli/internal/cmd/cmderror"
	"github.com/debricked/fortify-sca-mcp/v26/pkg/server"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestNewStartCmdMissingToken(t *testing.T) {
	token := ""
	cmd := NewStartCmd(&token, fakeAuthenticator{err: errors.New("not logged in")}, "https://debricked.com")
	cmd.SetArgs([]string{})

	err := cmd.Execute()

	assert.Error(t, err)
	var cmdErr cmderror.CommandError
	assert.ErrorAs(t, err, &cmdErr)
	assert.Contains(t, err.Error(), "not logged in")
}

func TestNewStartCmdExplicitTokenTakesPrecedence(t *testing.T) {
	original := serveFn
	defer func() { serveFn = original }()
	originalVerify := verifyFn
	defer func() { verifyFn = originalVerify }()
	verifyFn = func(context.Context, server.Options) error { return nil }

	var gotToken string
	serveFn = func(_ context.Context, options server.Options, _ io.Reader, _ io.Writer) error {
		gotToken = options.AccessToken

		return nil
	}

	token := "explicit-token"
	authenticator := fakeAuthenticator{token: &oauth2.Token{RefreshToken: "cached-refresh-token"}} //nolint:gosec // test fixture, not a real credential
	cmd := NewStartCmd(&token, authenticator, "https://debricked.com")

	err := cmd.Execute()

	assert.NoError(t, err)
	assert.Equal(t, "explicit-token", gotToken)
}

func TestNewStartCmdFallsBackToCachedLogin(t *testing.T) {
	original := serveFn
	defer func() { serveFn = original }()
	originalVerify := verifyFn
	defer func() { verifyFn = originalVerify }()
	verifyFn = func(context.Context, server.Options) error { return nil }

	var gotToken string
	serveFn = func(_ context.Context, options server.Options, _ io.Reader, _ io.Writer) error {
		gotToken = options.AccessToken

		return nil
	}

	token := ""
	authenticator := fakeAuthenticator{token: &oauth2.Token{ //nolint:gosec // test fixture, not a real credential
		AccessToken:  "cached-jwt",
		RefreshToken: "cached-refresh-token",
	}}
	cmd := NewStartCmd(&token, authenticator, "https://debricked.com")

	err := cmd.Execute()

	assert.NoError(t, err)
	assert.Equal(t, "cached-refresh-token", gotToken)
}

func TestNewStartCmdPropagatesBaseURL(t *testing.T) {
	original := serveFn
	defer func() { serveFn = original }()
	originalVerify := verifyFn
	defer func() { verifyFn = originalVerify }()

	var gotVerifyBaseURL, gotServeBaseURL, gotServeAPIVersion string
	verifyFn = func(_ context.Context, options server.Options) error {
		gotVerifyBaseURL = options.BaseURL

		return nil
	}
	serveFn = func(_ context.Context, options server.Options, _ io.Reader, _ io.Writer) error {
		gotServeBaseURL = options.BaseURL
		gotServeAPIVersion = options.APIVersion

		return nil
	}

	token := "token123"
	cmd := NewStartCmd(&token, fakeAuthenticator{}, "https://on-prem.example.com")

	err := cmd.Execute()

	assert.NoError(t, err)
	assert.Equal(t, "https://on-prem.example.com", gotVerifyBaseURL)
	assert.Equal(t, "https://on-prem.example.com", gotServeBaseURL)
	assert.Equal(t, "1.0", gotServeAPIVersion)
}

func TestNewStartCmdVerifyFailurePreventsServe(t *testing.T) {
	originalServe := serveFn
	defer func() { serveFn = originalServe }()
	originalVerify := verifyFn
	defer func() { verifyFn = originalVerify }()

	serveCalled := false
	serveFn = func(context.Context, server.Options, io.Reader, io.Writer) error {
		serveCalled = true

		return nil
	}
	verifyFn = func(context.Context, server.Options) error {
		return server.ErrUnauthorized
	}

	token := "bad-token"
	cmd := NewStartCmd(&token, fakeAuthenticator{}, "https://debricked.com")

	err := cmd.Execute()

	assert.Error(t, err)
	var cmdErr cmderror.CommandError
	assert.ErrorAs(t, err, &cmdErr)
	assert.False(t, serveCalled, "serveFn should not run when verification fails")
}

func TestNewStartCmdContextPropagation(t *testing.T) {
	original := serveFn
	defer func() { serveFn = original }()
	originalVerify := verifyFn
	defer func() { verifyFn = originalVerify }()
	verifyFn = func(context.Context, server.Options) error { return nil }

	var gotCtx context.Context
	serveFn = func(ctx context.Context, options server.Options, _ io.Reader, _ io.Writer) error {
		gotCtx = ctx
		assert.Equal(t, "token123", options.AccessToken)

		return nil
	}

	token := "token123"
	cmd := NewStartCmd(&token, fakeAuthenticator{}, "https://debricked.com")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd.SetContext(ctx)

	err := cmd.Execute()

	assert.NoError(t, err)
	assert.Equal(t, ctx, gotCtx)
}

func TestNewStartCmdContextCanceledIsCleanExit(t *testing.T) {
	original := serveFn
	defer func() { serveFn = original }()
	originalVerify := verifyFn
	defer func() { verifyFn = originalVerify }()
	verifyFn = func(context.Context, server.Options) error { return nil }

	serveFn = func(_ context.Context, _ server.Options, _ io.Reader, _ io.Writer) error {
		return context.Canceled
	}

	token := "token123"
	cmd := NewStartCmd(&token, fakeAuthenticator{}, "https://debricked.com")

	err := cmd.Execute()

	assert.NoError(t, err)
}

func TestNewStartCmdServeError(t *testing.T) {
	original := serveFn
	defer func() { serveFn = original }()
	originalVerify := verifyFn
	defer func() { verifyFn = originalVerify }()
	verifyFn = func(context.Context, server.Options) error { return nil }

	serveFn = func(_ context.Context, _ server.Options, _ io.Reader, _ io.Writer) error {
		return errors.New("boom")
	}

	token := "token123"
	cmd := NewStartCmd(&token, fakeAuthenticator{}, "https://debricked.com")

	err := cmd.Execute()

	assert.Error(t, err)
	var cmdErr cmderror.CommandError
	assert.ErrorAs(t, err, &cmdErr)
}

func TestNewStartCmdNoStdout(t *testing.T) {
	original := serveFn
	defer func() { serveFn = original }()
	originalVerify := verifyFn
	defer func() { verifyFn = originalVerify }()
	verifyFn = func(context.Context, server.Options) error { return nil }
	serveFn = func(_ context.Context, _ server.Options, _ io.Reader, _ io.Writer) error {
		return nil
	}

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	assert.NoError(t, err)
	os.Stdout = w
	defer func() { os.Stdout = oldStdout }()

	token := "token123"
	cmd := NewStartCmd(&token, fakeAuthenticator{}, "https://debricked.com")
	execErr := cmd.Execute()

	_ = w.Close()
	os.Stdout = oldStdout
	out, readErr := io.ReadAll(r)
	assert.NoError(t, readErr)

	assert.NoError(t, execErr)
	assert.Empty(t, out)
}

func TestResolveToken(t *testing.T) {
	tests := []struct {
		name          string
		explicit      string
		authenticator fakeAuthenticator
		wantToken     string
		wantErr       bool
	}{
		{
			name:          "explicit token wins over cached login",
			explicit:      "explicit-token",
			authenticator: fakeAuthenticator{token: &oauth2.Token{RefreshToken: "cached-refresh-token"}}, //nolint:gosec // test fixture, not a real credential
			wantToken:     "explicit-token",
		},
		{ //nolint:gosec // test fixture, not a real credential
			name:          "falls back to cached refresh token",
			explicit:      "",
			authenticator: fakeAuthenticator{token: &oauth2.Token{AccessToken: "jwt", RefreshToken: "cached-refresh-token"}}, //nolint:gosec // test fixture, not a real credential
			wantToken:     "cached-refresh-token",
		},
		{
			name:          "authenticator error surfaces as error",
			explicit:      "",
			authenticator: fakeAuthenticator{err: errors.New("keyring unavailable")},
			wantToken:     "",
			wantErr:       true,
		},
		{
			name:          "cached login with empty refresh token errors",
			explicit:      "",
			authenticator: fakeAuthenticator{token: &oauth2.Token{}},
			wantToken:     "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := resolveToken(tt.explicit, tt.authenticator)

			assert.Equal(t, tt.wantToken, gotToken)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
