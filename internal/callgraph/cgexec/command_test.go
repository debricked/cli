package cgexec

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	execTestdata "github.com/debricked/cli/internal/callgraph/cgexec/testdata"
)

func TestMakeCommandFailsWithContext(t *testing.T) {
	path, _ := exec.LookPath("mvn")
	ctx, cancel := NewContext(100)
	defer cancel()
	osCmd := MakeCommand(".", path, []string{
		"mvn",
		"package",
		"-q",
		"-DskipTests",
		"-e",
	}, ctx)
	cmd := NewCommand(osCmd)
	err := RunCommand(*cmd, ctx)
	t.Log(err)
	assert.Contains(t, err.Error(), "there is no POM in this directory")
}

func TestMakeCommandFailsWithNoContext(t *testing.T) {
	path, _ := exec.LookPath("mvn")
	osCmd := MakeCommand(".", path, []string{
		"mvn",
		"package",
		"-q",
		"-DskipTests",
		"-e",
	}, nil)
	cmd := NewCommand(osCmd)
	err := RunCommand(*cmd, nil)
	t.Log(err)
	assert.Contains(t, err.Error(), "exec: Stderr already set")
}

func TestMakeCommandDeadlineExceeded(t *testing.T) {
	ctx, _ := execTestdata.NewContextMockDeadlineReached()
	cmd := execTestdata.NewCommandMock()
	err := RunCommand(cmd, ctx)
	t.Log(err)
	if err == nil {
		assert.FailNow(t, "Error was unexpectedly nil.")
	}
	assert.Contains(t, err.Error(), "Timeout error: Set timeout duration for Callgraph jobs reached")
}

func TestMakeCommandCancelled(t *testing.T) {
	ctx, _ := execTestdata.NewContextMockCancelled()
	cmd := execTestdata.NewCommandMock()
	err := RunCommand(cmd, ctx)
	t.Log(err)
	assert.Contains(t, err.Error(), "Timeout error: Set timeout duration for Callgraph jobs reached")
}

func TestMakeCommandStartFailure(t *testing.T) {
	ctx, _ := execTestdata.NewContextMock()
	cmdConfig := execTestdata.NewCmdConfig()
	cmd := execTestdata.NewCommandMockWithConfig(*cmdConfig)
	err := RunCommand(cmd, ctx)
	t.Log(err)
	assert.Contains(t, err.Error(), "test error")
}

func TestMakeCommandSuccess(t *testing.T) {
	ctx, _ := execTestdata.NewContextMock()
	cmd := execTestdata.NewCommandMock()
	err := RunCommand(cmd, ctx)
	t.Log(err)
	assert.Nil(t, err)
}
