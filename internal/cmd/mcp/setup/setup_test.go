package setup

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/debricked/cli/internal/mcpclient"
	"github.com/debricked/cli/internal/mcpclient/testdata"
	"github.com/debricked/fortify-sca-mcp/v26/pkg/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func withStubbedVerify(t *testing.T, err error) {
	t.Helper()
	original := verifyFn
	verifyFn = func(context.Context, server.Options) error { return err }
	t.Cleanup(func() { verifyFn = original })
}

func withStubbedLookPath(t *testing.T, found bool) {
	t.Helper()
	original := lookPathFn
	if found {
		lookPathFn = func(string) (string, error) { return "/usr/bin/debricked", nil }
	} else {
		lookPathFn = func(string) (string, error) { return "", exec.ErrNotFound }
	}
	t.Cleanup(func() { lookPathFn = original })
}

func TestSetupCmdRequiresClientFlag(t *testing.T) {
	withStubbedVerify(t, nil)
	withStubbedLookPath(t, true)
	token := "tok"
	cmd := NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", testdata.NewFileSystem(), testdata.Environment{Wd: "/proj"})
	cmd.SetArgs([]string{})

	err := cmd.Execute()

	assert.Error(t, err)
}

func TestSetupCmdInvalidClient(t *testing.T) {
	withStubbedVerify(t, nil)
	withStubbedLookPath(t, true)
	token := "tok"
	cmd := NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", testdata.NewFileSystem(), testdata.Environment{Wd: "/proj"})
	cmd.SetArgs([]string{"--client", "bogus"})

	err := cmd.Execute()

	assert.Error(t, err)
}

func TestSetupCmdCreatesProjectConfig(t *testing.T) {
	withStubbedVerify(t, nil)
	withStubbedLookPath(t, true)
	fs := testdata.NewFileSystem()
	token := "tok"
	cmd := NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", fs, testdata.Environment{Wd: "/proj", GOOSValue: "linux"})
	cmd.SetArgs([]string{"--client", "vscode"})
	configPath := filepath.Join("/proj", ".vscode", "mcp.json")

	err := cmd.Execute()

	require.NoError(t, err)
	content, ok := fs.Files[configPath]
	require.True(t, ok)
	assert.Contains(t, string(content), `"debricked"`)
	assert.Contains(t, string(content), `"type": "stdio"`)
	assert.NotContains(t, fs.Files, configPath+".bak")
}

func TestSetupCmdFallsBackToAbsolutePathWhenNotOnPATH(t *testing.T) {
	withStubbedVerify(t, nil)
	withStubbedLookPath(t, false)
	fs := testdata.NewFileSystem()
	token := "tok"
	env := testdata.Environment{Wd: "/proj", GOOSValue: "linux", ExecutablePath: "/opt/debricked/debricked"}
	cmd := NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", fs, env)
	cmd.SetArgs([]string{"--client", "vscode"})
	configPath := filepath.Join("/proj", ".vscode", "mcp.json")

	err := cmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, string(fs.Files[configPath]), "/opt/debricked/debricked")
}

func TestSetupCmdIdempotent(t *testing.T) {
	withStubbedVerify(t, nil)
	withStubbedLookPath(t, true)
	fs := testdata.NewFileSystem()
	token := "tok"
	env := testdata.Environment{Wd: "/proj", GOOSValue: "linux"}
	cmd := NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", fs, env)
	cmd.SetArgs([]string{"--client", "vscode"})
	require.NoError(t, cmd.Execute())

	cmd2 := NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", fs, env)
	cmd2.SetArgs([]string{"--client", "vscode"})
	require.NoError(t, cmd2.Execute())

	assert.NotContains(t, fs.Files, filepath.Join("/proj", ".vscode", "mcp.json")+".bak")
	assert.NotContains(t, fs.Files, filepath.Join("/proj", ".github", "copilot-instructions.md")+".bak")
}

func TestSetupCmdInstallsProjectInstructionsForEveryClientAndScope(t *testing.T) {
	withStubbedVerify(t, nil)
	withStubbedLookPath(t, true)
	for _, tt := range []struct {
		client string
		path   string
	}{
		{"vscode", filepath.Join("/proj", ".github", "copilot-instructions.md")},
		{"claude", filepath.Join("/proj", "CLAUDE.md")},
		{"cursor", filepath.Join("/proj", ".cursor", "rules", "debricked.mdc")},
	} {
		for _, global := range []bool{false, true} {
			t.Run(tt.client+fmt.Sprintf("/global=%t", global), func(t *testing.T) {
				fs := testdata.NewFileSystem()
				token := "tok"
				env := testdata.Environment{Wd: "/proj", HomeDir: "/home/test", GOOSValue: "linux"}
				cmd := NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", fs, env)
				var stderr bytes.Buffer
				cmd.SetErr(&stderr)
				args := []string{"--client", tt.client}
				if global {
					args = append(args, "--global")
				}
				cmd.SetArgs(args)

				require.NoError(t, cmd.Execute())

				raw, ok := fs.Files[tt.path]
				require.True(t, ok)
				assert.Contains(t, string(raw), "## Dependency policy")
				assert.Contains(t, string(raw), "`check_dependency_policy_compliance`")
				assert.Contains(t, string(raw), "Get fresh approval before checking a different candidate list.")
				assert.NotContains(t, fs.Files, tt.path+mcpclient.BackupSuffix)
				if tt.client == "cursor" {
					assert.Contains(t, string(raw), "---\nalwaysApply: true\n---")
				}
				if tt.client == "claude" {
					assert.Contains(t, stderr.String(), "not automatically by Claude Desktop")
				}
			})
		}
	}
}

func TestSetupCmdRefreshesInstructionsWhenServerAlreadyConfigured(t *testing.T) {
	withStubbedVerify(t, nil)
	withStubbedLookPath(t, true)
	fs := testdata.NewFileSystem()
	token := "tok"
	env := testdata.Environment{Wd: "/proj", GOOSValue: "linux"}
	cmd := NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", fs, env)
	cmd.SetArgs([]string{"--client", "vscode"})
	require.NoError(t, cmd.Execute())
	path := filepath.Join("/proj", ".github", "copilot-instructions.md")
	original := "# My instructions\n\n<!-- debricked:dependency-policy:start -->\nOld policy\n<!-- debricked:dependency-policy:end -->\n\nKeep this.\n"
	fs.Files[path] = []byte(original)

	cmd = NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", fs, env)
	cmd.SetArgs([]string{"--client", "vscode"})
	require.NoError(t, cmd.Execute())

	assert.Contains(t, string(fs.Files[path]), "# My instructions")
	assert.Contains(t, string(fs.Files[path]), "Keep this.")
	assert.Contains(t, string(fs.Files[path]), "`check_dependency_policy_compliance`")
	assert.NotContains(t, string(fs.Files[path]), "Old policy")
	assert.NotContains(t, fs.Files, path+mcpclient.BackupSuffix)
	assert.NotContains(t, fs.Files, filepath.Join("/proj", ".vscode", "mcp.json")+mcpclient.BackupSuffix)

	delete(fs.Files, path)
	cmd = NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", fs, env)
	cmd.SetArgs([]string{"--client", "vscode"})
	require.NoError(t, cmd.Execute())
	assert.Contains(t, string(fs.Files[path]), "`check_dependency_policy_compliance`")
}

func TestSetupCmdReportsMalformedInstructions(t *testing.T) {
	withStubbedVerify(t, nil)
	withStubbedLookPath(t, true)
	fs := testdata.NewFileSystem()
	token := "tok"
	path := filepath.Join("/proj", ".github", "copilot-instructions.md")
	original := "# Keep this\n<!-- debricked:dependency-policy:start -->\n"
	fs.Files[path] = []byte(original)
	cmd := NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", fs, testdata.Environment{Wd: "/proj"})
	cmd.SetArgs([]string{"--client", "vscode"})

	err := cmd.Execute()

	require.ErrorContains(t, err, "invalid dependency policy markers")
	assert.Equal(t, original, string(fs.Files[path]))
	assert.NotContains(t, fs.Files, path+mcpclient.BackupSuffix)
}

func TestSetupCmdBacksUpExistingFile(t *testing.T) {
	withStubbedVerify(t, nil)
	withStubbedLookPath(t, true)
	fs := testdata.NewFileSystem()
	configPath := filepath.Join("/proj", ".vscode", "mcp.json")
	fs.Files[configPath] = []byte(`{"servers":{"other":{"command":"foo"}}}`)
	token := "tok"
	cmd := NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", fs, testdata.Environment{Wd: "/proj", GOOSValue: "linux"})
	cmd.SetArgs([]string{"--client", "vscode"})

	err := cmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, fs.Files, configPath+".bak")
	assert.Contains(t, string(fs.Files[configPath]), "other")
}

func TestSetupCmdMalformedExistingFileAborts(t *testing.T) {
	withStubbedVerify(t, nil)
	withStubbedLookPath(t, true)
	fs := testdata.NewFileSystem()
	fs.Files[filepath.Join("/proj", ".vscode", "mcp.json")] = []byte(`{ not json`)
	token := "tok"
	cmd := NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", fs, testdata.Environment{Wd: "/proj", GOOSValue: "linux"})
	cmd.SetArgs([]string{"--client", "vscode"})

	err := cmd.Execute()

	assert.Error(t, err)
}

func TestSetupCmdWarnsWhenNoCredentials(t *testing.T) {
	withStubbedVerify(t, nil)
	withStubbedLookPath(t, true)
	fs := testdata.NewFileSystem()
	token := ""
	cmd := NewSetupCmd(&token, fakeAuthenticator{err: errors.New("not logged in")}, "https://debricked.com", fs, testdata.Environment{Wd: "/proj", GOOSValue: "linux"})
	cmd.SetArgs([]string{"--client", "vscode"})

	err := cmd.Execute()

	// config is still written successfully even without credentials
	require.NoError(t, err)
	assert.Contains(t, fs.Files, filepath.Join("/proj", ".vscode", "mcp.json"))
}
