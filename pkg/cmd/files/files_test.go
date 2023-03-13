package files

import (
	"testing"

	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/client/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewFilesCmd(t *testing.T) {
	var debClient client.IDebClient = testdata.NewDebClientMock()
	cmd := NewFilesCmd(&debClient)
	commands := cmd.Commands()
	nbrOfCommands := 1
	assert.Lenf(t, commands, nbrOfCommands, "failed to assert that there were %d sub commands connected", nbrOfCommands)
}

func TestPreRun(t *testing.T) {
	var c client.IDebClient = testdata.NewDebClientMock()
	cmd := NewFilesCmd(&c)
	cmd.PreRun(cmd, nil)
}
