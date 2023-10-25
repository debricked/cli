package maven

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeDependencyTreeCmd(t *testing.T) {
	cmd, _ := CmdFactory{}.MakeDependencyTreeCmd(".")
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "mvn")
	assert.Contains(t, args, "dependency:tree")
	assert.Contains(t, args, "-DoutputFile=maven.debricked.lock")
	assert.Contains(t, args, "-DoutputType=tgf")
	assert.Contains(t, args, "--fail-at-end")
}
