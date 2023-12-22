package npm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeInstallCmd(t *testing.T) {
	npmCommand := "npm"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeInstallCmd(npmCommand, "file")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "npm")
	assert.Contains(t, args, "install")
}
