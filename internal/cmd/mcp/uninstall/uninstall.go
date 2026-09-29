package uninstall

import (
	"fmt"

	"github.com/debricked/cli/internal/cmd/cmderror"
	"github.com/debricked/cli/internal/io"
	"github.com/debricked/cli/internal/mcpclient"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// NewUninstallCmd creates the `mcp uninstall` command, which removes the debricked
// MCP server entry from a supported client's config file.
func NewUninstallCmd(fs io.IFileSystem, env mcpclient.IEnvironment) *cobra.Command {
	var clientFlag string
	var global bool

	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Remove the Debricked MCP server from a client's configuration.",
		Long: `Remove the "debricked" MCP server entry from a supported client's configuration
file (VS Code, Claude Desktop, or Cursor). The config file itself, and any other entries
in it, are left in place; a backup of the previous file is written to <file>.bak.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SilenceUsage = true

			client, err := mcpclient.ParseClient(clientFlag)
			if err != nil {
				return cmderror.CommandError{Code: 1, Err: err}
			}

			if global && client == mcpclient.ClientClaude {
				_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "note: Claude Desktop only has a global configuration; --global has no effect.")
			}

			path, err := mcpclient.ConfigPath(client, global, env)
			if err != nil {
				return cmderror.CommandError{Code: 1, Err: err}
			}

			config, existed, err := mcpclient.ReadConfig(fs, path)
			if err != nil {
				return cmderror.CommandError{Code: 1, Err: err}
			}

			if !existed || !mcpclient.RemoveServerEntry(config, mcpclient.ServerMapKey(client), mcpclient.ServerEntryName) {
				_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s Already not configured at %s.\n", color.GreenString("✔"), path)

				return nil
			}

			if err := mcpclient.WriteConfig(fs, path, config, existed); err != nil {
				return cmderror.CommandError{Code: 1, Err: err}
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Backed up previous config to %s%s\n", path, mcpclient.BackupSuffix)
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s Removed %s configuration from %s\n", color.GreenString("✔"), client, path)

			return nil
		},
	}

	cmd.Flags().StringVar(&clientFlag, "client", "", "MCP client to remove: vscode, claude, or cursor.")
	cmd.Flags().BoolVar(&global, "global", false, "Target the user-level (global) config instead of the project-level one.")
	_ = cmd.MarkFlagRequired("client")

	return cmd
}
