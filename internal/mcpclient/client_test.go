package mcpclient_test

import (
	"path/filepath"
	"testing"

	"github.com/debricked/cli/internal/mcpclient"
	"github.com/debricked/cli/internal/mcpclient/testdata"
	"github.com/stretchr/testify/assert"
)

func TestParseClient(t *testing.T) {
	tests := []struct {
		input   string
		want    mcpclient.Client
		wantErr bool
	}{
		{"vscode", mcpclient.ClientVSCode, false},
		{"VSCode", mcpclient.ClientVSCode, false},
		{" cursor ", mcpclient.ClientCursor, false},
		{"CLAUDE", mcpclient.ClientClaude, false},
		{"windsurf", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		got, err := mcpclient.ParseClient(tt.input)
		if tt.wantErr {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		}
	}
}

func TestServerMapKey(t *testing.T) {
	assert.Equal(t, "servers", mcpclient.ServerMapKey(mcpclient.ClientVSCode))
	assert.Equal(t, "mcpServers", mcpclient.ServerMapKey(mcpclient.ClientClaude))
	assert.Equal(t, "mcpServers", mcpclient.ServerMapKey(mcpclient.ClientCursor))
}

func TestSupportsLocalScope(t *testing.T) {
	assert.True(t, mcpclient.SupportsLocalScope(mcpclient.ClientVSCode))
	assert.True(t, mcpclient.SupportsLocalScope(mcpclient.ClientCursor))
	assert.False(t, mcpclient.SupportsLocalScope(mcpclient.ClientClaude))
}

func TestConfigPathLocal(t *testing.T) {
	env := testdata.Environment{Wd: "/home/user/project", GOOSValue: "linux"}

	path, err := mcpclient.ConfigPath(mcpclient.ClientVSCode, false, env)
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join("/home/user/project", ".vscode", "mcp.json"), path)

	path, err = mcpclient.ConfigPath(mcpclient.ClientCursor, false, env)
	assert.NoError(t, err)
	assert.Equal(t, filepath.Join("/home/user/project", ".cursor", "mcp.json"), path)
}

func TestConfigPathClaudeAlwaysGlobal(t *testing.T) {
	env := testdata.Environment{Wd: "/home/user/project", GOOSValue: "darwin", HomeDir: "/Users/test"}

	// global=false must still resolve to the global path for claude.
	path, err := mcpclient.ConfigPath(mcpclient.ClientClaude, false, env)
	assert.NoError(t, err)
	assert.Equal(t, "/Users/test/Library/Application Support/Claude/claude_desktop_config.json", path)
}

func TestConfigPathGlobalPerOS(t *testing.T) {
	tests := []struct {
		name   string
		client mcpclient.Client
		env    testdata.Environment
		want   string
	}{
		{
			name:   "vscode mac",
			client: mcpclient.ClientVSCode,
			env:    testdata.Environment{GOOSValue: "darwin", HomeDir: "/Users/test"},
			want:   "/Users/test/Library/Application Support/Code/User/mcp.json",
		},
		{
			name:   "vscode windows",
			client: mcpclient.ClientVSCode,
			env:    testdata.Environment{GOOSValue: "windows", Env: map[string]string{"APPDATA": `C:\Users\test\AppData\Roaming`}},
			want:   `C:\Users\test\AppData\Roaming\Code\User\mcp.json`,
		},
		{
			name:   "vscode linux with XDG_CONFIG_HOME",
			client: mcpclient.ClientVSCode,
			env:    testdata.Environment{GOOSValue: "linux", Env: map[string]string{"XDG_CONFIG_HOME": "/home/test/.config"}},
			want:   "/home/test/.config/Code/User/mcp.json",
		},
		{
			name:   "vscode linux fallback to home",
			client: mcpclient.ClientVSCode,
			env:    testdata.Environment{GOOSValue: "linux", HomeDir: "/home/test"},
			want:   "/home/test/.config/Code/User/mcp.json",
		},
		{
			name:   "claude mac",
			client: mcpclient.ClientClaude,
			env:    testdata.Environment{GOOSValue: "darwin", HomeDir: "/Users/test"},
			want:   "/Users/test/Library/Application Support/Claude/claude_desktop_config.json",
		},
		{
			name:   "claude windows",
			client: mcpclient.ClientClaude,
			env:    testdata.Environment{GOOSValue: "windows", Env: map[string]string{"APPDATA": `C:\Users\test\AppData\Roaming`}},
			want:   `C:\Users\test\AppData\Roaming\Claude\claude_desktop_config.json`,
		},
		{
			name:   "claude linux best-effort",
			client: mcpclient.ClientClaude,
			env:    testdata.Environment{GOOSValue: "linux", HomeDir: "/home/test"},
			want:   "/home/test/.config/Claude/claude_desktop_config.json",
		},
		{
			name:   "cursor mac",
			client: mcpclient.ClientCursor,
			env:    testdata.Environment{GOOSValue: "darwin", HomeDir: "/Users/test"},
			want:   "/Users/test/.cursor/mcp.json",
		},
		{
			name:   "cursor windows",
			client: mcpclient.ClientCursor,
			env:    testdata.Environment{GOOSValue: "windows", HomeDir: `C:\Users\test`},
			want:   `C:\Users\test\.cursor\mcp.json`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, err := mcpclient.ConfigPath(tt.client, true, tt.env)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, path)
		})
	}
}

func TestConfigPathMissingAppData(t *testing.T) {
	env := testdata.Environment{GOOSValue: "windows"}

	_, err := mcpclient.ConfigPath(mcpclient.ClientVSCode, true, env)
	assert.Error(t, err)

	_, err = mcpclient.ConfigPath(mcpclient.ClientClaude, true, env)
	assert.Error(t, err)
}

func TestJoinPathWindows(t *testing.T) {
	assert.Equal(t, `a\b\c`, mcpclient.JoinPath("windows", "a", "b", "c"))
	assert.Equal(t, "a/b/c", mcpclient.JoinPath("linux", "a", "b", "c"))
}
