package yarn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeExecPath struct{}

func (fakeExecPath) LookPath(file string) (string, error) {
	// Simulate a successful lookup of the yarn binary without requiring it in PATH.
	return file, nil
}

func TestMakeInstallCmd(t *testing.T) {
	yarnCommand := "yarn"
	cmd, err := CmdFactory{
		execPath: fakeExecPath{},
	}.MakeInstallCmd(yarnCommand, "file")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "yarn")
	assert.Contains(t, args, "install")
}
