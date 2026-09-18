package mcp

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

// apiVersion matches the Debricked API version the rest of the CLI targets (e.g. /api/1.0/...).
const apiVersion = "1.0"

// NewStartCmd creates the `mcp start` command, which runs an MCP server over stdio.
func NewStartCmd(accessToken *string, authenticator auth.IAuthenticator, baseURL string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the Debricked MCP server over stdio.",
		Long:  `Start a Model Context Protocol (MCP) server, communicating over stdin/stdout, until the client disconnects or the process is interrupted.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SilenceUsage = true

			token, resolveErr := resolveToken(*accessToken, authenticator)
			if token == "" {
				err := errors.New("no access token found. Pass --access-token/DEBRICKED_TOKEN, or run `debricked auth login` first")
				if resolveErr != nil {
					err = fmt.Errorf("%w (%v)", err, resolveErr)
				}

				return cmderror.CommandError{Code: 1, Err: err}
			}

			options := server.Options{AccessToken: token, BaseURL: baseURL, APIVersion: apiVersion}

			// Fail fast on bad credentials instead of only discovering it on the first tool call.
			if err := verifyFn(cmd.Context(), options); err != nil {
				return cmderror.CommandError{Code: 1, Err: fmt.Errorf("access token rejected: %w", err)}
			}

			// stdout is reserved for MCP JSON-RPC traffic, so status goes to stderr.
			fmt.Fprintln(cmd.ErrOrStderr(), "debricked MCP server authenticated, running on stdio. Hit Ctrl-C to exit.")

			err := serveFn(cmd.Context(), options, os.Stdin, os.Stdout)
			if err == nil || errors.Is(err, context.Canceled) {
				return nil
			}

			return cmderror.CommandError{Code: 1, Err: err}
		},
	}

	return cmd
}

// resolveToken prefers an explicit access token, falling back to the CLI's cached login.
// The Fortify SCA MCP server treats this value as a refresh token, so the cached
// session's RefreshToken is used rather than its short-lived JWT AccessToken.
// The returned error (only set when the token is empty) carries the underlying
// cause, e.g. a keyring failure vs. simply never having logged in.
func resolveToken(accessToken string, authenticator auth.IAuthenticator) (string, error) {
	if token := strings.TrimSpace(accessToken); token != "" {
		return token, nil
	}

	cachedToken, err := authenticator.Token()
	if err != nil {
		return "", err
	}
	if cachedToken == nil || strings.TrimSpace(cachedToken.RefreshToken) == "" {
		return "", errors.New("cached login has no refresh token")
	}

	return strings.TrimSpace(cachedToken.RefreshToken), nil
}
