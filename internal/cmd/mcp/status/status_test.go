package status

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/debricked/cli/internal/mcpclient/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusCmdInvalidClient(t *testing.T) {
	cmd := NewStatusCmd(testdata.NewFileSystem(), testdata.Environment{Wd: "/proj", GOOSValue: "linux", HomeDir: "/home/test"})
	cmd.SetArgs([]string{"--client", "bogus"})

	err := cmd.Execute()

	assert.Error(t, err)
}

func TestStatusCmdNotConfiguredAnywhere(t *testing.T) {
	var out bytes.Buffer
	cmd := NewStatusCmd(testdata.NewFileSystem(), testdata.Environment{Wd: "/proj", GOOSValue: "linux", HomeDir: "/home/test"})
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--client", "claude"})

	err := cmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, out.String(), "not configured")
}

func TestStatusCmdConfiguredAndWorking(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-based handshake fixtures are not supported on windows")
	}
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not available")
	}

	fs := testdata.NewFileSystem()
	configPath := filepath.Join("/home/test", ".config", "Claude", "claude_desktop_config.json")
	fs.Files[configPath] = []byte(
		`{"mcpServers":{"debricked":{"command":"sh","args":["-c","read line; echo '{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":{}}'"]}}}`,
	)
	var out bytes.Buffer
	cmd := NewStatusCmd(fs, testdata.Environment{GOOSValue: "linux", HomeDir: "/home/test"})
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--client", "claude"})

	err := cmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, out.String(), "responded to a live MCP handshake")
}

func TestStatusCmdConfiguredButFailingHandshake(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-based handshake fixtures are not supported on windows")
	}
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("sh not available")
	}

	fs := testdata.NewFileSystem()
	configPath := filepath.Join("/home/test", ".config", "Claude", "claude_desktop_config.json")
	fs.Files[configPath] = []byte(
		`{"mcpServers":{"debricked":{"command":"sh","args":["-c","echo \"no access token found\" 1>&2; exit 1"]}}}`,
	)
	var out bytes.Buffer
	cmd := NewStatusCmd(fs, testdata.Environment{GOOSValue: "linux", HomeDir: "/home/test"})
	cmd.SetOut(&out)
	cmd.SetArgs([]string{"--client", "claude"})

	err := cmd.Execute()

	assert.Error(t, err)
	assert.Contains(t, out.String(), "server check failed")
}
