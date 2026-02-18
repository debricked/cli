package pnpm

import (
	"testing"

	"github.com/debricked/cli/internal/resolution/pm/pnpm/testdata"
	"github.com/stretchr/testify/assert"
)

func TestMakeInstallCmd(t *testing.T) {
	pnpmCommand := "pnpm"
	cmd, err := CmdFactory{
		execPath: testdata.ExecPathMock{},
	}.MakeInstallCmd(pnpmCommand, "file")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "pnpm")
	assert.Contains(t, args, "install")
}
