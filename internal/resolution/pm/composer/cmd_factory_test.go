package composer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeExecPath struct{}

func (fakeExecPath) LookPath(file string) (string, error) {
	// Simulate a successful lookup of the composer binary without requiring it in PATH.
	return file, nil
}

func TestMakeInstallCmd(t *testing.T) {
	composerCommand := "composer"
	cmd, err := CmdFactory{
		execPath: fakeExecPath{},
	}.MakeInstallCmd(composerCommand, "file")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "composer")
	assert.Contains(t, args, "update")
}
