package pub

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type execPathMock struct{}

func (execPathMock) LookPath(file string) (string, error) {
	return "/usr/bin/" + file, nil
}

func TestMakeLockCmd(t *testing.T) {
	factory := CmdFactory{execPath: execPathMock{}}
	manifest := filepath.Join("some", "path", "pubspec.yaml")

	cmd, err := factory.MakeLockCmd(manifest)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Equal(t, "/usr/bin/dart", cmd.Path)
	assert.Contains(t, cmd.Args, "dart")
	assert.Contains(t, cmd.Args, "pub")
	assert.Contains(t, cmd.Args, "get")
	assert.Equal(t, filepath.Dir(manifest), cmd.Dir)
}

func TestMakeDepsCmd(t *testing.T) {
	factory := CmdFactory{execPath: execPathMock{}}
	manifest := filepath.Join("some", "path", "pubspec.yaml")

	cmd, err := factory.MakeDepsCmd(manifest)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Equal(t, "/usr/bin/dart", cmd.Path)
	assert.Contains(t, cmd.Args, "dart")
	assert.Contains(t, cmd.Args, "pub")
	assert.Contains(t, cmd.Args, "deps")
	assert.Contains(t, cmd.Args, "--json")
	assert.Equal(t, filepath.Dir(manifest), cmd.Dir)
}
