package pub

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/pub/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestRunCmdErrExecutableNotFound(t *testing.T) {
	execErr := errors.New("exec: \"dart\": executable file not found in $PATH")
	j := NewJob("file", testdata.CmdFactoryMock{LockErr: execErr})

	go jobTestdata.WaitStatus(j)

	j.Run()

	errs := j.Errors().GetAll()
	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "executable file not found")
	assert.Contains(t, errs[0].Documentation(), "Dart wasn't found")
}

func TestRunDepsCmdErrExecutableNotFound(t *testing.T) {
	execErr := errors.New("exec: \"dart\": executable file not found in $PATH")
	j := NewJob("file", testdata.CmdFactoryMock{Name: "echo", Arg: "ok", DepsErr: execErr})

	go jobTestdata.WaitStatus(j)

	j.Run()

	errs := j.Errors().GetAll()
	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "executable file not found")
	assert.Contains(t, errs[0].Documentation(), "Dart wasn't found")
}

func TestRunSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	manifest := filepath.Join(tmpDir, "pubspec.yaml")

	err := os.WriteFile(manifest, []byte("name: test"), 0600)
	assert.NoError(t, err)

	j := NewJob(manifest, testdata.CmdFactoryMock{Name: "echo", Arg: "ok"})

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.False(t, j.Errors().HasError())
	assert.Len(t, j.Errors().GetAll(), 0)

	_, statErr := os.Stat(filepath.Join(tmpDir, "pubspec.deps.json"))
	assert.NoError(t, statErr)
}
