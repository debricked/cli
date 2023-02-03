package gomod

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeGraphCmd(t *testing.T) {
	cmd, _ := CmdFactory{}.MakeGraphCmd()
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "go")
	assert.Contains(t, args, "mod")
	assert.Contains(t, args, "graph")
}

func TestMakeListCmd(t *testing.T) {
	cmd, _ := CmdFactory{}.MakeListCmd()
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "go")
	assert.Contains(t, args, "list")
	assert.Contains(t, args, "-mod=readonly")
	assert.Contains(t, args, "-e")
	assert.Contains(t, args, "-m")
	assert.Contains(t, args, "all")
}
