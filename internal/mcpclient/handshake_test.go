package mcpclient_test

import (
	"context"
	"os/exec"
	"runtime"
	"testing"
	"time"

	"github.com/debricked/cli/internal/mcpclient"
	"github.com/stretchr/testify/assert"
)

func requireShell(t *testing.T) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("shell-based handshake fixtures are not supported on windows")
	}
	shell, err := exec.LookPath("sh")
	if err != nil {
		t.Skip("sh not available")
	}

	return shell
}

func TestCheckServerSuccess(t *testing.T) {
	shell := requireShell(t)

	script := `read line; echo '{"jsonrpc":"2.0","id":1,"result":{}}'`
	err := mcpclient.CheckServer(context.Background(), shell, []string{"-c", script}, time.Second)

	assert.NoError(t, err)
}

func TestCheckServerJSONRPCError(t *testing.T) {
	shell := requireShell(t)

	script := `read line; echo '{"jsonrpc":"2.0","id":1,"error":{"code":-1,"message":"boom"}}'`
	err := mcpclient.CheckServer(context.Background(), shell, []string{"-c", script}, time.Second)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "boom")
}

func TestCheckServerAuthRejected(t *testing.T) {
	shell := requireShell(t)

	script := `echo "access token rejected: bad token" 1>&2; exit 1`
	err := mcpclient.CheckServer(context.Background(), shell, []string{"-c", script}, time.Second)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "access token rejected")
}

func TestCheckServerNoAccessToken(t *testing.T) {
	shell := requireShell(t)

	script := `echo "no access token found. Pass --access-token" 1>&2; exit 1`
	err := mcpclient.CheckServer(context.Background(), shell, []string{"-c", script}, time.Second)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no access token found")
}

func TestCheckServerTimeout(t *testing.T) {
	shell := requireShell(t)

	// "exec" replaces the shell process itself, so killing the spawned pid actually
	// stops the sleep instead of leaving an orphaned grandchild holding stdout open.
	script := `exec sleep 5`
	err := mcpclient.CheckServer(context.Background(), shell, []string{"-c", script}, 100*time.Millisecond)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "timed out")
}

func TestCheckServerCommandNotFound(t *testing.T) {
	err := mcpclient.CheckServer(context.Background(), "debricked-does-not-exist-binary", nil, time.Second)

	assert.Error(t, err)
}
