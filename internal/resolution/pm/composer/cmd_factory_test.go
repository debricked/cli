package composer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeInstallCmd(t *testing.T) {
	composerCommand := "composer"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeInstallCmd(composerCommand, "file")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "composer")
	assert.Contains(t, args, "update")
}
