package cgexec

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	execTestdata "github.com/debricked/cli/internal/callgraph/cgexec/testdata"
)

func TestMakeCommandWithContext(t *testing.T) {
	cmdName := "echo"
	path, _ := exec.LookPath(cmdName)
	ctx, cancel := NewContext(100)
	defer cancel()
	osCmd := MakeCommand(".", path, []string{
		"echo",
		"package",
	}, ctx)
	cmd := NewCommand(osCmd)
	err := RunCommand(*cmd, ctx)
	assert.Nil(t, err)
}

func TestMakeCommandWithNoContext(t *testing.T) {
	cmdName := "echo"
	path, _ := exec.LookPath(cmdName)
	osCmd := MakeCommand(".", path, []string{
		"echo",
		"package",
		"-q",
		"-DskipTests",
		"-e",
	}, nil)
	cmd := NewCommand(osCmd)
	err := RunCommand(*cmd, nil)
	assert.Nil(t, err)
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

func TestMakeCommandCancelledInteruptedError(t *testing.T) {
	ctx, _ := execTestdata.NewContextMockCancelled()
	cmd := execTestdata.NewCommandMock()
	cmd.Process = &os.Process{}
	cmd.SignalError = fmt.Errorf("error")
	err := RunCommand(cmd, ctx)
	t.Log(err)
	assert.Contains(t, err.Error(), "error")
}

func TestMakeCommandStartFailure(t *testing.T) {
	ctx, _ := execTestdata.NewContextMock()
	cmdConfig := execTestdata.NewCmdConfig()
	cmd := execTestdata.NewCommandMockWithConfig(*cmdConfig)
	err := RunCommand(cmd, ctx)
	t.Log(err)
	assert.Contains(t, err.Error(), "test error")
}

func TestWaitExecutionFail(t *testing.T) {
	ctx, _ := execTestdata.NewContextMock()
	cmd := execTestdata.NewCommandMock()
	cmd.WaitError = fmt.Errorf("error")
	err := RunCommand(cmd, ctx)
	t.Log(err)
	assert.Contains(t, err.Error(), "error")
}

func TestNoContextCommandFail(t *testing.T) {
	cmd := execTestdata.NewCommandMock()
	cmd.CombinedOutputError = &exec.ExitError{}
	err := RunCommand(cmd, nil)
	t.Log(err)
	assert.Contains(t, err.Error(), "executed in folder")
}

func TestMakeCommandSuccess(t *testing.T) {
	ctx, _ := execTestdata.NewContextMock()
	cmd := execTestdata.NewCommandMock()
	err := RunCommand(cmd, ctx)
	t.Log(err)
	assert.Nil(t, err)
}

func TestCmdGetProcess(t *testing.T) {
	cmd := NewCommand(exec.Command("echo", "GetProcess"))
	process := cmd.GetProcess()

	assert.Nil(t, process)
}

func TestCmdGetDir(t *testing.T) {
	cmd := NewCommand(exec.Command("echo", "GetProcess"))
	dir := cmd.GetDir()

	assert.NotNil(t, dir)
}
