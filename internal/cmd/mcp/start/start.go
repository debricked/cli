package start

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/debricked/cli/internal/auth"
	"github.com/debricked/cli/internal/cmd/cmderror"
	"github.com/debricked/fortify-sca-mcp/v26/pkg/server"
	"github.com/spf13/cobra"
)

// serveFn is a seam over server.Serve to allow tests to substitute a fake implementation.
var serveFn = server.Serve

// verifyFn is a seam over server.VerifyAccessToken to allow tests to substitute a fake implementation.
var verifyFn = server.VerifyAccessToken

// APIVersion matches the Debricked API version the rest of the CLI targets (e.g. /api/1.0/...).
const APIVersion = "1.0"

// NewStartCmd creates the `mcp start` command, which runs an MCP server over stdio.
func NewStartCmd(accessToken *string, authenticator auth.IAuthenticator, baseURL string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the Debricked MCP server over stdio.",
		Long:  `Start a Model Context Protocol (MCP) server, communicating over stdin/stdout, until the client disconnects or the process is interrupted.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SilenceUsage = true

			options, err := buildServerOptions(*accessToken, authenticator, baseURL)
			if err != nil {
				return cmderror.CommandError{Code: 1, Err: err}
			}

			if err := verifyFn(cmd.Context(), options); err != nil {
				if options.Authenticate == nil {
					return cmderror.CommandError{Code: 1, Err: fmt.Errorf("access token rejected: %w", err)}
				}
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "debricked MCP login required; use the authenticate tool to sign in.")
			}

			// stdout is reserved for MCP JSON-RPC traffic, so status goes to stderr.
			if _, err := fmt.Fprintln(cmd.ErrOrStderr(), "debricked MCP server running on stdio. API authentication is checked on tool calls. Hit Ctrl-C to exit."); err != nil {
				return cmderror.CommandError{Code: 1, Err: err}
			}

			err = serveFn(cmd.Context(), options, os.Stdin, os.Stdout)
			if err == nil || errors.Is(err, context.Canceled) {
				return nil
			}

			return cmderror.CommandError{Code: 1, Err: err}
		},
	}

	return cmd
}

// buildServerOptions decides how the MCP server authenticates. An explicit
// --access-token/DEBRICKED_TOKEN (a PAT) is passed through as-is: the server exchanges
// it itself via /api/login_refresh, exactly as before. Otherwise this falls back to the
// cached `debricked auth login` session via a TokenFetcher callback, so the server
// always gets a currently-valid bearer JWT on demand - the OAuth refresh token itself
// is never handed over, since /api/login_refresh doesn't accept it.
func buildServerOptions(explicitToken string, authenticator auth.IAuthenticator, baseURL string) (server.Options, error) {
	if token := strings.TrimSpace(explicitToken); token != "" {
		return server.Options{AccessToken: token, BaseURL: baseURL, APIVersion: APIVersion}, nil
	}

	return server.Options{
		TokenFetcher: auth.NewCachedTokenFetcher(authenticator),
		Authenticate: func(ctx context.Context) error {
			if contextual, ok := authenticator.(interface {
				AuthenticateContext(context.Context) error
			}); ok {
				return contextual.AuthenticateContext(ctx)
			}
			return authenticator.Authenticate()
		},
		BaseURL:    baseURL,
		APIVersion: APIVersion,
	}, nil
}
