package yarn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeInstallCmd(t *testing.T) {
	yarnCommand := "yarn"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeInstallCmd(yarnCommand, "file")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "yarn")
	assert.Contains(t, args, "install")
}
