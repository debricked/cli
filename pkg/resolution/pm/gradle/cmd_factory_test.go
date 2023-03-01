package gradle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeDependenciesCmd(t *testing.T) {
	cmd, _ := CmdFactory{}.MakeDependenciesCmd(".")
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "gradle")
	assert.Contains(t, args, "dependencies")
}
