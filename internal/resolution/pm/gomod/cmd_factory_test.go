package gomod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeGraphCmd(t *testing.T) {
	cmd, _ := CmdFactory{}.MakeGraphCmd(".")
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "go")
	assert.Contains(t, args, "mod")
	assert.Contains(t, args, "graph")
}

func TestMakeListCmd(t *testing.T) {
	cmd, _ := CmdFactory{}.MakeListCmd(".")
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "go")
	assert.Contains(t, args, "list")
	assert.Contains(t, args, "-mod=readonly")
	assert.Contains(t, args, "-e")
	assert.Contains(t, args, "-m")
	assert.Contains(t, args, "all")
}

func TestMakeListJsonCmd(t *testing.T) {
	factory := CmdFactory{}
	cmd, err := factory.MakeListJsonCmd(".")
	assert.Nil(t, err)
	assert.NotNil(t, cmd)
	assert.Contains(t, cmd.Args, "list")
	assert.Contains(t, cmd.Args, "-json")
	assert.Contains(t, cmd.Args, "./...")
}
