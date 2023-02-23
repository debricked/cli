package pip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateVenvCmd(t *testing.T) {
	venvName := "test-file.venv"
	cmd, _ := CmdFactory{}.MakeCreateVenvCmd(venvName)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "python")
	assert.Contains(t, args, "-m")
	assert.Contains(t, args, "venv")
	assert.Contains(t, args, venvName)
	assert.Contains(t, args, "--clear")
}

func TestMakeInstallCmd(t *testing.T) {
	fileName := "test-file"
	pipCommand := "pip"
	cmd, _ := CmdFactory{}.MakeInstallCmd(pipCommand, fileName)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "pip")
	assert.Contains(t, args, "install")
	assert.Contains(t, args, "-r")
	assert.Contains(t, args, fileName)
}

func TestMakeCatCmd(t *testing.T) {
	fileName := "test-file"
	cmd, _ := CmdFactory{}.MakeCatCmd(fileName)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "cat")
	assert.Contains(t, args, fileName)
}
func TestMakeListCmd(t *testing.T) {
	mockCommand := "mock-cmd"
	cmd, _ := CmdFactory{}.MakeListCmd(mockCommand)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "pip")
	assert.Contains(t, args, "list")
}

func TestMakeShowCmd(t *testing.T) {
	input := []string{"package1", "package2"}
	mockCommand := "pip"
	cmd, _ := CmdFactory{}.MakeShowCmd(mockCommand, input)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "pip")
	assert.Contains(t, args, "show")
	assert.Contains(t, args, "package1")
	assert.Contains(t, args, "package2")
}
