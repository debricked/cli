package setup

import (
	"context"
	"errors"
	"os/exec"
	"testing"

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

	err := cmd.Execute()

	require.NoError(t, err)
	content, ok := fs.Files["/proj/.vscode/mcp.json"]
	require.True(t, ok)
	assert.Contains(t, string(content), `"debricked"`)
	assert.NotContains(t, fs.Files, "/proj/.vscode/mcp.json.bak")
}

func TestSetupCmdFallsBackToAbsolutePathWhenNotOnPATH(t *testing.T) {
	withStubbedVerify(t, nil)
	withStubbedLookPath(t, false)
	fs := testdata.NewFileSystem()
	token := "tok"
	env := testdata.Environment{Wd: "/proj", GOOSValue: "linux", ExecutablePath: "/opt/debricked/debricked"}
	cmd := NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", fs, env)
	cmd.SetArgs([]string{"--client", "vscode"})

	err := cmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, string(fs.Files["/proj/.vscode/mcp.json"]), "/opt/debricked/debricked")
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

	assert.NotContains(t, fs.Files, "/proj/.vscode/mcp.json.bak")
}

func TestSetupCmdBacksUpExistingFile(t *testing.T) {
	withStubbedVerify(t, nil)
	withStubbedLookPath(t, true)
	fs := testdata.NewFileSystem()
	fs.Files["/proj/.vscode/mcp.json"] = []byte(`{"servers":{"other":{"command":"foo"}}}`)
	token := "tok"
	cmd := NewSetupCmd(&token, fakeAuthenticator{}, "https://debricked.com", fs, testdata.Environment{Wd: "/proj", GOOSValue: "linux"})
	cmd.SetArgs([]string{"--client", "vscode"})

	err := cmd.Execute()

	require.NoError(t, err)
	assert.Contains(t, fs.Files, "/proj/.vscode/mcp.json.bak")
	assert.Contains(t, string(fs.Files["/proj/.vscode/mcp.json"]), "other")
}

func TestSetupCmdMalformedExistingFileAborts(t *testing.T) {
	withStubbedVerify(t, nil)
	withStubbedLookPath(t, true)
	fs := testdata.NewFileSystem()
	fs.Files["/proj/.vscode/mcp.json"] = []byte(`{ not json`)
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
	assert.Contains(t, fs.Files, "/proj/.vscode/mcp.json")
}
