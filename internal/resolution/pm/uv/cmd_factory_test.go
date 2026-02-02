package uv

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
	manifest := filepath.Join("testdata", "pyproject.toml")

	cmd, err := factory.MakeLockCmd(manifest)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Equal(t, "/usr/bin/uv", cmd.Path)
	assert.Contains(t, cmd.Args, "uv")
	assert.Contains(t, cmd.Args, "lock")
	assert.Equal(t, filepath.Dir(manifest), cmd.Dir)
}
