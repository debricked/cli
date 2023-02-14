package pip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeListCmd(t *testing.T) {
	cmd, _ := CmdFactory{}.MakeListCmd()
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "pip")
	assert.Contains(t, args, "list")
}

func TestMakeShowCmd(t *testing.T) {
	input := []string{"package1", "package2"}
	cmd, _ := CmdFactory{}.MakeShowCmd(input)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "pip")
	assert.Contains(t, args, "show")
	assert.Contains(t, args, "package1")
	assert.Contains(t, args, "package2")
}
