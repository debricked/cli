package mcpclient

import (
	"fmt"
	"path/filepath"
	"strings"

	internalOs "github.com/debricked/cli/internal/runtime/os"
)

// Client identifies a supported MCP host application.
type Client string

const (
	ClientVSCode Client = "vscode"
	ClientClaude Client = "claude"
	ClientCursor Client = "cursor"

	// ServerEntryName is the fixed key used for the debricked entry in a client's
	// server map, e.g. servers.debricked / mcpServers.debricked.
	ServerEntryName = "debricked"
)

// AllClients lists every supported client, in a stable order used by `mcp status`
// when no --client flag is given.
var AllClients = []Client{ClientVSCode, ClientClaude, ClientCursor}

const (
	macOS = "darwin"
	winOS = "windows"
)

// ParseClient validates a --client flag value, case-insensitively.
func ParseClient(value string) (Client, error) {
	switch Client(strings.ToLower(strings.TrimSpace(value))) {
	case ClientVSCode:
		return ClientVSCode, nil
	case ClientClaude:
		return ClientClaude, nil
	case ClientCursor:
		return ClientCursor, nil
	default:
		return "", fmt.Errorf("unsupported client %q, must be one of: vscode, claude, cursor", value)
	}
}

// ServerMapKey returns the top-level JSON key holding the server map for a given client.
func ServerMapKey(client Client) string {
	if client == ClientVSCode {
		return "servers"
	}

	return "mcpServers" // claude, cursor
}

// SupportsLocalScope reports whether a client has a project-level config location.
// Claude Desktop only has a global config.
func SupportsLocalScope(client Client) bool {
	return client != ClientClaude
}

// ConfigPath resolves the config file path for a client/scope combination.
func ConfigPath(client Client, global bool, env IEnvironment) (string, error) {
	if client == ClientClaude || global {
		return globalConfigPath(client, env)
	}

	return localConfigPath(client, env)
}

func localConfigPath(client Client, env IEnvironment) (string, error) {
	wd, err := env.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to resolve current directory: %w", err)
	}

	var dir string
	switch client {
	case ClientVSCode:
		dir = ".vscode"
	case ClientCursor:
		dir = ".cursor"
	default:
		return "", fmt.Errorf("client %q has no project-level config", client)
	}

	return filepath.Join(wd, dir, "mcp.json"), nil
}

func globalConfigPath(client Client, env IEnvironment) (string, error) {
	switch client {
	case ClientClaude:
		return claudeGlobalPath(env)
	case ClientVSCode:
		return vscodeGlobalPath(env)
	case ClientCursor:
		home, err := env.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory: %w", err)
		}

		return JoinPath(env.GOOS(), home, ".cursor", "mcp.json"), nil
	default:
		return "", fmt.Errorf("unsupported client %q", client)
	}
}

func claudeGlobalPath(env IEnvironment) (string, error) {
	goos := env.GOOS()
	switch goos {
	case macOS:
		home, err := env.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory: %w", err)
		}

		return JoinPath(goos, home, "Library", "Application Support", "Claude", "claude_desktop_config.json"), nil
	case winOS:
		appData := env.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("%%APPDATA%% is not set")
		}

		return JoinPath(goos, appData, "Claude", "claude_desktop_config.json"), nil
	default:
		// Claude Desktop doesn't officially ship on Linux; this mirrors the community-maintained location.
		home, err := env.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory: %w", err)
		}

		return JoinPath(goos, home, ".config", "Claude", "claude_desktop_config.json"), nil
	}
}

func vscodeGlobalPath(env IEnvironment) (string, error) {
	goos := env.GOOS()
	switch goos {
	case macOS:
		home, err := env.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory: %w", err)
		}

		return JoinPath(goos, home, "Library", "Application Support", "Code", "User", "mcp.json"), nil
	case winOS:
		appData := env.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("%%APPDATA%% is not set")
		}

		return JoinPath(goos, appData, "Code", "User", "mcp.json"), nil
	default:
		if xdg := env.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return JoinPath(goos, xdg, "Code", "User", "mcp.json"), nil
		}
		home, err := env.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to resolve home directory: %w", err)
		}

		return JoinPath(goos, home, ".config", "Code", "User", "mcp.json"), nil
	}
}

// JoinPath joins path segments using the separator appropriate for goos, so path
// construction can be tested for a given target OS regardless of the host test OS.
func JoinPath(goos string, parts ...string) string {
	sep := "/"
	if goos == internalOs.Windows {
		sep = "\\"
	}

	return strings.Join(parts, sep)
}
