package status

import (
	"errors"
	"fmt"

	"github.com/debricked/cli/internal/cmd/cmderror"
	"github.com/debricked/cli/internal/io"
	"github.com/debricked/cli/internal/mcpclient"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewStatusCmd creates the `mcp status` command, which reports whether the debricked
// MCP server is configured for one or all supported clients, and whether it actually
// starts and authenticates.
func NewStatusCmd(fs io.IFileSystem, env mcpclient.IEnvironment) *cobra.Command {
	var clientFlag string

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show whether the Debricked MCP server is configured and working.",
		Long: `Check whether the "debricked" MCP server entry is present in supported clients'
configuration files, and perform a live MCP handshake against the exact command found
there to confirm it starts and authenticates successfully. Checks vscode, claude, and
cursor if --client is omitted.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SilenceUsage = true

			clients := mcpclient.AllClients
			if clientFlag != "" {
				client, err := mcpclient.ParseClient(clientFlag)
				if err != nil {
					return cmderror.CommandError{Code: 1, Err: err}
				}
				clients = []mcpclient.Client{client}
			}

			ok := true
			for _, client := range clients {
				for _, global := range scopesFor(client) {
					if !checkClientScope(cmd, fs, env, client, global) {
						ok = false
					}
				}
			}

			if !ok {
				return cmderror.CommandError{Code: 1, Err: errors.New("one or more checks failed, see above")}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&clientFlag, "client", "", "Only check one client: vscode, claude, or cursor.")

	return cmd
}

func scopesFor(client mcpclient.Client) []bool {
	if !mcpclient.SupportsLocalScope(client) {
		return []bool{true}
	}

	return []bool{false, true}
}

func scopeLabel(global bool) string {
	if global {
		return "global"
	}

	return "project"
}

// checkClientScope reports the config + live-handshake status for one client/scope
// combination and returns false if anything is missing or failing.
func checkClientScope(cmd *cobra.Command, fs io.IFileSystem, env mcpclient.IEnvironment, client mcpclient.Client, global bool) bool {
	out := cmd.OutOrStdout()
	label := fmt.Sprintf("%s (%s)", client, scopeLabel(global))

	path, err := mcpclient.ConfigPath(client, global, env)
	if err != nil {
		fmt.Fprintf(out, "%s %s: %v\n", color.RedString("✘"), label, err)

		return false
	}

	config, existed, err := mcpclient.ReadConfig(fs, path)
	if err != nil {
		fmt.Fprintf(out, "%s %s: %v\n", color.RedString("✘"), label, err)

		return false
	}

	entry, ok := mcpclient.HasServerEntry(config, mcpclient.ServerMapKey(client), mcpclient.ServerEntryName)
	if !existed || !ok {
		fmt.Fprintf(out, "%s %s: not configured (%s)\n", color.RedString("✘"), label, path)

		return false
	}

	command, args, err := mcpclient.EntryCommand(entry)
	if err != nil {
		fmt.Fprintf(out, "%s %s: configured at %s but entry is invalid: %v\n", color.RedString("✘"), label, path, err)

		return false
	}
	fmt.Fprintf(out, "%s %s: configured (%s)\n", color.GreenString("✔"), label, path)

	if err := mcpclient.CheckServer(cmd.Context(), command, args, mcpclient.DefaultHandshakeTimeout); err != nil {
		fmt.Fprintf(out, "  %s server check failed: %v\n", color.RedString("✘"), err)

		return false
	}
	fmt.Fprintf(out, "  %s server responded to a live MCP handshake\n", color.GreenString("✔"))

	return true
}
