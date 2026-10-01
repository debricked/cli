package uninstall

import (
	"path/filepath"
	"testing"

	"github.com/debricked/cli/internal/mcpclient/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUninstallCmdRequiresClientFlag(t *testing.T) {
	cmd := NewUninstallCmd(testdata.NewFileSystem(), testdata.Environment{Wd: "/proj"})
	cmd.SetArgs([]string{})

	err := cmd.Execute()

	assert.Error(t, err)
}

func TestUninstallCmdAlreadyNotConfigured(t *testing.T) {
	fs := testdata.NewFileSystem()
	cmd := NewUninstallCmd(fs, testdata.Environment{Wd: "/proj", GOOSValue: "linux"})
	configPath := filepath.Join("/proj", ".vscode", "mcp.json")
	cmd.SetArgs([]string{"--client", "vscode"})

	err := cmd.Execute()

	require.NoError(t, err)
	assert.NotContains(t, fs.Files, configPath)
}

func TestUninstallCmdRemovesEntryKeepsFile(t *testing.T) {
	fs := testdata.NewFileSystem()
	configPath := filepath.Join("/proj", ".vscode", "mcp.json")
	fs.Files[configPath] = []byte(`{"servers":{"debricked":{"command":"debricked","args":["mcp","start"]},"other":{"command":"foo"}}}`)
	cmd := NewUninstallCmd(fs, testdata.Environment{Wd: "/proj", GOOSValue: "linux"})
	cmd.SetArgs([]string{"--client", "vscode"})

	err := cmd.Execute()

	require.NoError(t, err)
	content := string(fs.Files[configPath])
	assert.NotContains(t, content, `"debricked":{`)
	assert.Contains(t, content, "other")
	assert.Contains(t, fs.Files, configPath+".bak")
}

func TestUninstallCmdMalformedExistingFileAborts(t *testing.T) {
	fs := testdata.NewFileSystem()
	fs.Files[filepath.Join("/proj", ".vscode", "mcp.json")] = []byte(`{ not json`)
	cmd := NewUninstallCmd(fs, testdata.Environment{Wd: "/proj", GOOSValue: "linux"})
	cmd.SetArgs([]string{"--client", "vscode"})

	err := cmd.Execute()

	assert.Error(t, err)
}
