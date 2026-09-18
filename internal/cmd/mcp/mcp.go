package mcp

import (
	"github.com/debricked/cli/internal/auth"
	"github.com/spf13/cobra"
)

// NewMCPCmd creates the parent `mcp` command, grouping MCP related subcommands.
// accessToken is the resolved root-level access token (from --access-token / DEBRICKED_TOKEN).
// authenticator provides a fallback to the CLI's cached login when accessToken is empty.
// baseURL is the Debricked host (respects DEBRICKED_URI) the MCP server should talk to.
func NewMCPCmd(accessToken *string, authenticator auth.IAuthenticator, baseURL string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Debricked MCP server.",
		Long:  `Run a Model Context Protocol (MCP) server exposing Debricked functionality.`,
	}
	cmd.AddCommand(NewStartCmd(accessToken, authenticator, baseURL))

	return cmd
}
