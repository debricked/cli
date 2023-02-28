package files

import (
	"testing"

	"github.com/debricked/cli/pkg/file"
	"github.com/debricked/cli/pkg/resolution"
	"github.com/stretchr/testify/assert"
)

func TestNewFilesCmd(t *testing.T) {
	finder, _ := file.NewFinder(nil)
	resolver := resolution.NewResolver(
		nil,
		nil,
		nil,
	)
	cmd := NewFilesCmd(finder, resolver)
	commands := cmd.Commands()
	nbrOfCommands := 2
	assert.Lenf(t, commands, nbrOfCommands, "failed to assert that there were %d sub commands connected", nbrOfCommands)
}

func TestPreRun(t *testing.T) {
	var c client.IDebClient = testdata.NewDebClientMock()
	cmd := NewFilesCmd(&c)
	cmd.PreRun(cmd, nil)
}
