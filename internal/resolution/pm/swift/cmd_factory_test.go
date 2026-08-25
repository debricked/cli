package swift

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type execPathMock struct{}

type failingExecPathMock struct{}

func (execPathMock) LookPath(file string) (string, error) {
	return "/usr/bin/" + file, nil
}

func (failingExecPathMock) LookPath(_ string) (string, error) {
	return "", errors.New("swift executable not found")
}

func TestMakeResolveCmd(t *testing.T) {
	factory := CmdFactory{execPath: execPathMock{}}
	manifest := filepath.Join("some", "path", "Package.swift")

	cmd, err := factory.MakeResolveCmd(manifest)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Equal(t, "/usr/bin/swift", cmd.Path)
	assert.Contains(t, cmd.Args, "swift")
	assert.Contains(t, cmd.Args, "package")
	assert.Contains(t, cmd.Args, "resolve")
	assert.Equal(t, filepath.Dir(manifest), cmd.Dir)
}

func TestMakeDepsCmd(t *testing.T) {
	factory := CmdFactory{execPath: execPathMock{}}
	manifest := filepath.Join("some", "path", "Package.swift")

	cmd, err := factory.MakeDepsCmd(manifest)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	assert.Equal(t, "/usr/bin/swift", cmd.Path)
	assert.Contains(t, cmd.Args, "swift")
	assert.Contains(t, cmd.Args, "package")
	assert.Contains(t, cmd.Args, "show-dependencies")
	assert.Contains(t, cmd.Args, "--format")
	assert.Contains(t, cmd.Args, "json")
	assert.Equal(t, filepath.Dir(manifest), cmd.Dir)
}

func TestMakeResolveCmdLookupError(t *testing.T) {
	factory := CmdFactory{execPath: failingExecPathMock{}}

	cmd, err := factory.MakeResolveCmd("Package.swift")

	assert.Nil(t, cmd)
	assert.EqualError(t, err, "swift executable not found")
}

func TestMakeDepsCmdLookupError(t *testing.T) {
	factory := CmdFactory{execPath: failingExecPathMock{}}

	cmd, err := factory.MakeDepsCmd("Package.swift")

	assert.Nil(t, cmd)
	assert.EqualError(t, err, "swift executable not found")
}
