package nuget

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeInstallCmd(t *testing.T) {
	nugetCommand := "dotnet"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeInstallCmd(nugetCommand, "file")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "dotnet")
	assert.Contains(t, args, "restore")
}
