package files

import (
	"testing"

	"github.com/debricked/cli/internal/file"
	"github.com/stretchr/testify/assert"
)

func TestNewFilesCmd(t *testing.T) {
	finder, _ := file.NewFinder(nil)
	cmd := NewFilesCmd(finder, nil)
	commands := cmd.Commands()
	nbrOfCommands := 2
	assert.Lenf(t, commands, nbrOfCommands, "failed to assert that there were %d sub commands connected", nbrOfCommands)
}

func TestPreRun(t *testing.T) {
	cmd := NewFilesCmd(nil, nil)
	cmd.PreRun(cmd, nil)
}
