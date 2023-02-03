package maven

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeDependencyTreeCmd(t *testing.T) {
	cmd, _ := CmdFactory{}.MakeDependencyTreeCmd()
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "mvn")
	assert.Contains(t, args, "dependency:tree")
	assert.Contains(t, args, "-DoutputFile=.debricked-maven-dependencies.tgf")
	assert.Contains(t, args, "-DoutputType=tgf")
}
