package cgexec

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeCommandFailsWithContext(t *testing.T) {
	path, _ := exec.LookPath("mvn")
	ctx, cancel := NewContext(100)
	defer cancel()
	cmd := MakeCommand(".", path, []string{
		"mvn",
		"package",
		"-q",
		"-DskipTests",
		"-e",
	}, ctx)
	err := RunCommand(cmd, ctx)
	t.Log(err)
	assert.Contains(t, err.Error(), "there is no POM in this directory")
}

func TestMakeCommandFailsWithNoContext(t *testing.T) {
	path, _ := exec.LookPath("mvn")
	cmd := MakeCommand(".", path, []string{
		"mvn",
		"package",
		"-q",
		"-DskipTests",
		"-e",
	}, nil)
	err := RunCommand(cmd, nil)
	t.Log(err)
	assert.Contains(t, err.Error(), "exec: Stderr already set")
}

// func TestMakeCommandFailsDeadlineExceeded(t *testing.T) {
// 	ctx, _ := execTestdata.NewContextMock()
// 	cmd := execTestdata.NewCmdMock()
// 	err := RunCommand(cmd, ctx)
// 	t.Log(err)
// 	assert.Contains(t, err.Error(), "Timeout error: Set timeout duration for Callgraph jobs reached")
// }