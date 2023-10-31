package pip

import (
	"errors"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ExecPathMock struct {
	python3Error error
	pythonError  error
}

func (epm ExecPathMock) LookPath(file string) (string, error) {
	if epm.python3Error != nil && file == "python3" {
		return "", epm.python3Error
	}

	if epm.pythonError != nil && file == "python" {
		return "", epm.pythonError
	}

	return file, nil
}

func TestCreateVenvCmd(t *testing.T) {
	venvName := "test-file.venv"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeCreateVenvCmd(venvName)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "python3")
	assert.Contains(t, args, "-m")
	assert.Contains(t, args, "venv")
	assert.Contains(t, args, venvName)
	assert.Contains(t, args, "--clear")
}

func TestCreateVenvCmdPython3Error(t *testing.T) {
	err := errors.New("executable file not found in $PATH")
	execPathMock := ExecPathMock{python3Error: err}
	venvName := "test-file-python3-error.venv"
	cmd, err := CmdFactory{
		execPath: execPathMock,
	}.MakeCreateVenvCmd(venvName)

	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "python")
	assert.Contains(t, args, "-m")
	assert.Contains(t, args, "venv")
	assert.Contains(t, args, venvName)
	assert.Contains(t, args, "--clear")
}

func TestCreateVenvCmdPythonCompletelyMissing(t *testing.T) {
	pathErr := errors.New("executable file not found in $PATH")
	execPathMock := ExecPathMock{python3Error: pathErr, pythonError: pathErr}
	venvName := "test-file-python-missing.venv"
	_, err := CmdFactory{
		execPath: execPathMock,
	}.MakeCreateVenvCmd(venvName)

	assert.ErrorContains(t, err, "executable file not found in")
	assert.ErrorContains(t, err, "PATH")
}

func TestMakeInstallCmd(t *testing.T) {
	fileName := "test-file"
	pipCommand := "pip"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeInstallCmd(pipCommand, fileName)
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "pip")
	assert.Contains(t, args, "install")
	assert.Contains(t, args, "-r")
	assert.Contains(t, args, fileName)
}

func TestMakeCatCmd(t *testing.T) {
	fileName := "test-file"
	expectedCommand := "cat"
	if runtime.GOOS == "windows" {
		expectedCommand = "type"
	}
	cmd, _ := CmdFactory{
		execPath: ExecPath{},
	}.MakeCatCmd(fileName)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, expectedCommand)
	assert.Contains(t, args, fileName)
}
func TestMakeListCmd(t *testing.T) {
	mockCommand := "mock-cmd"
	cmd, _ := CmdFactory{
		execPath: ExecPath{},
	}.MakeListCmd(mockCommand)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "pip")
	assert.Contains(t, args, "list")
}

func TestMakeShowCmd(t *testing.T) {
	input := []string{"package1", "package2"}
	mockCommand := "pip"
	cmd, _ := CmdFactory{
		execPath: ExecPath{},
	}.MakeShowCmd(mockCommand, input)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "pip")
	assert.Contains(t, args, "show")
	assert.Contains(t, args, "package1")
	assert.Contains(t, args, "package2")
}
