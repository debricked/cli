package bower

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeInstallCmd(t *testing.T) {
	bowerCommand := "bower"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeInstallCmd(bowerCommand, "file")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "bower")
	assert.Contains(t, args, "install")
}

func TestMakeListCmd(t *testing.T) {
	bowerCommand := "bower"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeListCmd(bowerCommand, "file")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "bower")
	assert.Contains(t, args, "list")
}
