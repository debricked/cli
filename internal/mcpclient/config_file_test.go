package mcpclient_test

import (
	"testing"

	"github.com/debricked/cli/internal/mcpclient"
	"github.com/debricked/cli/internal/mcpclient/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadConfigMissingFile(t *testing.T) {
	fs := testdata.NewFileSystem()

	config, existed, err := mcpclient.ReadConfig(fs, "/tmp/mcp.json")

	require.NoError(t, err)
	assert.False(t, existed)
	assert.Empty(t, config)
}

func TestReadConfigMalformedJSON(t *testing.T) {
	fs := testdata.NewFileSystem()
	fs.Files["/tmp/mcp.json"] = []byte(`{ not valid json`)

	_, _, err := mcpclient.ReadConfig(fs, "/tmp/mcp.json")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "/tmp/mcp.json")
}

func TestReadConfigPreservesUnrelatedKeys(t *testing.T) {
	fs := testdata.NewFileSystem()
	fs.Files["/tmp/mcp.json"] = []byte(`{"servers":{"other":{"command":"foo"}},"inputs":[1,2,3]}`)

	config, existed, err := mcpclient.ReadConfig(fs, "/tmp/mcp.json")

	require.NoError(t, err)
	assert.True(t, existed)
	assert.Contains(t, config, "inputs")
	assert.Contains(t, config, "servers")
}

func TestWriteConfigCreatesFileNoBackupWhenNew(t *testing.T) {
	fs := testdata.NewFileSystem()
	path := "/tmp/proj/.vscode/mcp.json"

	err := mcpclient.WriteConfig(fs, path, map[string]any{"servers": map[string]any{}}, false)

	require.NoError(t, err)
	assert.Contains(t, fs.Files, path)
	assert.NotContains(t, fs.Files, path+mcpclient.BackupSuffix)
}

func TestWriteConfigBacksUpExistingFile(t *testing.T) {
	fs := testdata.NewFileSystem()
	path := "/tmp/proj/.vscode/mcp.json"
	fs.Files[path] = []byte(`{"servers":{}}`)

	err := mcpclient.WriteConfig(fs, path, map[string]any{"servers": map[string]any{"debricked": map[string]any{}}}, true)

	require.NoError(t, err)
	assert.Equal(t, []byte(`{"servers":{}}`), fs.Files[path+mcpclient.BackupSuffix])
	assert.NotEqual(t, []byte(`{"servers":{}}`), fs.Files[path])
}

func TestUpsertServerEntryCreatesEntry(t *testing.T) {
	config := map[string]any{}

	changed := mcpclient.UpsertServerEntry(config, "servers", "debricked", map[string]any{
		"command": "debricked",
		"args":    []string{"mcp", "start"},
	})

	assert.True(t, changed)
	entry, ok := mcpclient.HasServerEntry(config, "servers", "debricked")
	require.True(t, ok)
	assert.Equal(t, "debricked", entry["command"])
}

func TestUpsertServerEntryPreservesOtherServersAndFields(t *testing.T) {
	config := map[string]any{
		"servers": map[string]any{
			"other": map[string]any{"command": "other-cmd"},
			"debricked": map[string]any{
				"command": "old-path",
				"args":    []string{"mcp", "start"},
				"env":     map[string]any{"FOO": "bar"},
			},
		},
	}

	changed := mcpclient.UpsertServerEntry(config, "servers", "debricked", map[string]any{
		"command": "debricked",
		"args":    []string{"mcp", "start"},
	})

	assert.True(t, changed)
	servers := config["servers"].(map[string]any)
	assert.Contains(t, servers, "other")
	entry := servers["debricked"].(map[string]any)
	assert.Equal(t, "debricked", entry["command"])
	assert.Equal(t, map[string]any{"FOO": "bar"}, entry["env"])
}

func TestUpsertServerEntryIdempotent(t *testing.T) {
	config := map[string]any{}
	fields := map[string]any{"command": "debricked", "args": []string{"mcp", "start"}}

	changed := mcpclient.UpsertServerEntry(config, "servers", "debricked", fields)
	assert.True(t, changed)

	changed = mcpclient.UpsertServerEntry(config, "servers", "debricked", fields)
	assert.False(t, changed)
}

func TestRemoveServerEntry(t *testing.T) {
	config := map[string]any{
		"mcpServers": map[string]any{
			"debricked": map[string]any{"command": "debricked"},
			"other":     map[string]any{"command": "other"},
		},
	}

	removed := mcpclient.RemoveServerEntry(config, "mcpServers", "debricked")
	assert.True(t, removed)

	_, ok := mcpclient.HasServerEntry(config, "mcpServers", "debricked")
	assert.False(t, ok)
	_, ok = mcpclient.HasServerEntry(config, "mcpServers", "other")
	assert.True(t, ok)
}

func TestRemoveServerEntryMissing(t *testing.T) {
	config := map[string]any{}

	removed := mcpclient.RemoveServerEntry(config, "mcpServers", "debricked")
	assert.False(t, removed)
}

func TestEntryCommandInvalid(t *testing.T) {
	_, _, err := mcpclient.EntryCommand(map[string]any{})
	assert.Error(t, err)

	_, _, err = mcpclient.EntryCommand(map[string]any{"command": "debricked", "args": []any{1, 2}})
	assert.Error(t, err)
}
