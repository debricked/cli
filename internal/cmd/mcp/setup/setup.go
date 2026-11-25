package setup

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/debricked/cli/internal/auth"
	"github.com/debricked/cli/internal/cmd/cmderror"
	"github.com/debricked/cli/internal/io"
	"github.com/debricked/cli/internal/mcpclient"
	"github.com/debricked/fortify-sca-mcp/v26/pkg/server"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// apiVersion matches the Debricked API version the rest of the CLI targets (e.g. /api/1.0/...).
const apiVersion = "1.0"

// verifyFn is a seam over server.VerifyAccessToken to allow tests to substitute a fake implementation.
var verifyFn = server.VerifyAccessToken

// lookPathFn is a seam over exec.LookPath to allow tests to substitute a fake implementation.
var lookPathFn = exec.LookPath

// NewSetupCmd creates the `mcp setup` command, which adds a debricked MCP server
// entry to a supported client's config file and installs project instructions.
func NewSetupCmd(accessToken *string, authenticator auth.IAuthenticator, baseURL string, fs io.IFileSystem, env mcpclient.IEnvironment) *cobra.Command {
	var clientFlag string
	var global bool

	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Configure the Debricked MCP server and project dependency policy instructions.",
		Long: `Add a "debricked" MCP server entry to a supported client's configuration file
(VS Code, Claude Desktop, or Cursor), so it can be picked up without hand-editing JSON.

By default, VS Code and Cursor are configured at the project level (./.vscode/mcp.json,
./.cursor/mcp.json, relative to the current directory). Pass --global to configure the
user-level config instead. Claude Desktop only has a global configuration; --global has
no effect for it.

Dependency policy instructions are always installed in the current project:
VS Code: .github/copilot-instructions.md
Claude: CLAUDE.md
Cursor: .cursor/rules/debricked.mdc (always applied)
Existing instructions are preserved; rerunning setup refreshes the managed policy block.
CLAUDE.md is read by Claude Code, not automatically by Claude Desktop; for Desktop,
copy the policy into your project's instructions manually.

No access token is ever written to the config file. Run 'debricked auth login' beforehand
so the server can authenticate when your MCP client starts it.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runSetup(cmd, clientFlag, global, accessToken, authenticator, baseURL, fs, env)
		},
	}

	cmd.Flags().StringVar(&clientFlag, "client", "", "MCP client to configure: vscode, claude, or cursor.")
	cmd.Flags().BoolVar(&global, "global", false, "Configure the user-level (global) config instead of the project-level one.")
	_ = cmd.MarkFlagRequired("client")

	return cmd
}

func runSetup(cmd *cobra.Command, clientFlag string, global bool, accessToken *string, authenticator auth.IAuthenticator, baseURL string, fs io.IFileSystem, env mcpclient.IEnvironment) error {
	cmd.SilenceUsage = true

	client, err := mcpclient.ParseClient(clientFlag)
	if err != nil {
		return cmderror.CommandError{Code: 1, Err: err}
	}

	if global && client == mcpclient.ClientClaude {
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "note: Claude Desktop only has a global configuration; --global has no effect.")
	}

	path, instructionsPath, err := resolveSetupPaths(client, global, env)
	if err != nil {
		return cmderror.CommandError{Code: 1, Err: err}
	}

	changed, err := configureServer(cmd, client, path, fs, env)
	if err != nil {
		return cmderror.CommandError{Code: 1, Err: err}
	}

	if err := installInstructions(cmd, client, instructionsPath, fs); err != nil {
		return cmderror.CommandError{Code: 1, Err: err}
	}
	if client == mcpclient.ClientClaude {
		_, _ = fmt.Fprintln(cmd.ErrOrStderr(), "note: CLAUDE.md is read by Claude Code, not automatically by Claude Desktop; copy the policy into your Desktop project's instructions manually.")
	}
	if changed {
		warnIfNoCredentials(cmd, accessToken, authenticator, baseURL)
	}

	return nil
}

func resolveSetupPaths(client mcpclient.Client, global bool, env mcpclient.IEnvironment) (configPath, instructionsPath string, err error) {
	configPath, err = mcpclient.ConfigPath(client, global, env)
	if err != nil {
		return "", "", err
	}

	instructionsPath, err = mcpclient.InstructionsPath(client, env)

	return configPath, instructionsPath, err
}

func configureServer(cmd *cobra.Command, client mcpclient.Client, path string, fs io.IFileSystem, env mcpclient.IEnvironment) (bool, error) {
	config, existed, err := mcpclient.ReadConfig(fs, path)
	if err != nil {
		return false, err
	}

	command, portable := resolveCommand(env)
	fields := map[string]any{
		"type":    "stdio",
		"command": command,
		"args":    []string{"mcp", "start"},
	}

	changed := mcpclient.UpsertServerEntry(config, mcpclient.ServerMapKey(client), mcpclient.ServerEntryName, fields)
	if !changed {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s MCP server already configured at %s, no config changes made.\n", color.GreenString("✔"), path)

		return false, nil
	}
	if err := mcpclient.WriteConfig(fs, path, config, existed); err != nil {
		return false, err
	}
	if existed {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "Backed up previous config to %s%s\n", path, mcpclient.BackupSuffix)
	}
	if !portable {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "note: %q was not found on PATH, using the current binary's absolute path instead. This is machine-specific and won't survive a reinstall or move.\n", "debricked")
	}
	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s Configured %s at %s\n", color.GreenString("✔"), client, path)

	return true, nil
}

func installInstructions(cmd *cobra.Command, client mcpclient.Client, path string, fs io.IFileSystem) error {
	changed, _, err := mcpclient.SetupInstructions(fs, path, client)
	if err != nil {
		return err
	}
	if changed {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s Installed dependency policy instructions at %s\n", color.GreenString("✔"), path)
	} else {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s Dependency policy instructions already up to date at %s\n", color.GreenString("✔"), path)
	}

	return nil
}

// resolveCommand prefers the bare "debricked" command, which resolves via the MCP
// client's own PATH and is safe to commit to a shared project config. It falls back
// to the absolute path of the currently running binary if "debricked" isn't on PATH.
// The second return value reports whether the portable form was used.
func resolveCommand(env mcpclient.IEnvironment) (string, bool) {
	if _, err := lookPathFn("debricked"); err == nil {
		return "debricked", true
	}

	if path, err := env.Executable(); err == nil {
		return path, false
	}

	return "debricked", true
}

// warnIfNoCredentials prints a best-effort warning (never an error) if no access
// token can currently be resolved, or if the resolved token is rejected by the API -
// the config was still written successfully either way.
func warnIfNoCredentials(cmd *cobra.Command, accessToken *string, authenticator auth.IAuthenticator, baseURL string) {
	var options server.Options
	if token := strings.TrimSpace(*accessToken); token != "" {
		options = server.Options{AccessToken: token, BaseURL: baseURL, APIVersion: apiVersion}
	} else if _, err := authenticator.Token(); err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "warning: no access token found; run `debricked auth login` before starting the server. (%v)\n", err)

		return
	} else {
		options = server.Options{TokenFetcher: auth.NewCachedTokenFetcher(authenticator), BaseURL: baseURL, APIVersion: apiVersion}
	}

	if err := verifyFn(cmd.Context(), options); err != nil {
		_, _ = fmt.Fprintf(cmd.ErrOrStderr(), "warning: access token rejected: %v\n", err)
	}
}
