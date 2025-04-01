package sbt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakePomCmd(t *testing.T) {
	cmd, _ := CmdFactory{}.MakePomCmd(".")
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "sbt")
	assert.Contains(t, args, "makePom")
}
