package bower

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeInstallCmd(t *testing.T) {
	bowerCommand := "bower"
	cmd, _ := CmdFactory{
		execPath: ExecPath{},
	}.MakeInstallCmd(bowerCommand, "file")
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "bower")
	assert.Contains(t, args, "install")
	assert.Contains(t, args, "--save")
	assert.Contains(t, args, "--save-dev")
	assert.Contains(t, args, "--save-exact")
	assert.Contains(t, args, "--allow-root")
}

func TestMakeListCmd(t *testing.T) {
	bowerCommand := "bower"
	cmd, _ := CmdFactory{
		execPath: ExecPath{},
	}.MakeListCmd(bowerCommand, "file")
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "bower")
	assert.Contains(t, args, "list")
}
