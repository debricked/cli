package swift

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/swift/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	j := NewJob("Package.swift", testdata.CmdFactoryMock{})
	assert.Equal(t, "Package.swift", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestRunCmdErrExecutableNotFound(t *testing.T) {
	execErr := errors.New("exec: \"swift\": executable file not found in $PATH")
	j := NewJob("Package.swift", testdata.CmdFactoryMock{LockErr: execErr})

	go jobTestdata.WaitStatus(j)
	j.Run()

	errs := j.Errors().GetAll()
	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "executable file not found")
	assert.Contains(t, errs[0].Documentation(), "Swift wasn't found")
}

func TestRunDepsCmdErrExecutableNotFound(t *testing.T) {
	execErr := errors.New("exec: \"swift\": executable file not found in $PATH")
	j := NewJob("Package.swift", testdata.CmdFactoryMock{Name: "echo", Arg: "ok", DepsErr: execErr})

	go jobTestdata.WaitStatus(j)
	j.Run()

	errs := j.Errors().GetAll()
	assert.Len(t, errs, 1)
	assert.Contains(t, errs[0].Error(), "executable file not found")
	assert.Contains(t, errs[0].Documentation(), "Swift wasn't found")
}

func TestRunSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	manifest := filepath.Join(tmpDir, "Package.swift")
	assert.NoError(t, os.WriteFile(manifest, []byte("// swift-tools-version: 5.9\n"), 0600))

	depsFile, err := filepath.Abs(filepath.Join("testdata", "dependencies.json"))
	assert.NoError(t, err)

	j := NewJob(manifest, testdata.CmdFactoryMock{Name: "echo", Arg: "ok", DepsFile: depsFile})
	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.False(t, j.Errors().HasError())
	lockContent, statErr := os.ReadFile(filepath.Join(tmpDir, ".spm.debricked.lock"))
	assert.NoError(t, statErr)

	var root dependencyNode
	assert.NoError(t, json.Unmarshal(lockContent, &root))
	assert.Equal(t, "example-app", root.Identity)
	assert.Len(t, root.Dependencies, 2)
	assert.Equal(t, "swift-nio", root.Dependencies[0].Identity)
	assert.Equal(t, "2.65.0", root.Dependencies[0].Version)
	assert.Len(t, root.Dependencies[0].Dependencies, 1)
	assert.Equal(t, "swift-collections", root.Dependencies[0].Dependencies[0].Identity)
}

func TestRunInvalidDependencyTree(t *testing.T) {
	tmpDir := t.TempDir()
	manifest := filepath.Join(tmpDir, "Package.swift")
	assert.NoError(t, os.WriteFile(manifest, []byte("// swift-tools-version: 5.9\n"), 0600))

	j := NewJob(manifest, testdata.CmdFactoryMock{Name: "echo", Arg: "ok"})
	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.True(t, j.Errors().HasError())
	assert.Contains(t, j.Errors().GetAll()[0].Error(), "did not return a JSON dependency tree")
	_, statErr := os.Stat(filepath.Join(tmpDir, ".spm.debricked.lock"))
	assert.Error(t, statErr)
}

func TestExtractDependencyTree(t *testing.T) {
	cases := []struct {
		name    string
		output  string
		wantErr string
	}{
		{name: "plain tree", output: `{"identity":"app","name":"App","version":"unspecified","dependencies":[]}`},
		{name: "tree with surrounding command noise", output: "Fetching package\n{\"identity\":\"app\",\"name\":\"App\",\"dependencies\":[]}\n"},
		{name: "no json", output: "error: could not resolve dependencies", wantErr: "did not return a JSON dependency tree"},
		{name: "malformed json", output: `{"identity": "app",}`, wantErr: "invalid character"},
		{name: "missing root package", output: `{"dependencies":[]}`, wantErr: "missing root package information"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			tree, err := extractDependencyTree([]byte(c.output))
			if c.wantErr != "" {
				assert.ErrorContains(t, err, c.wantErr)
				assert.Nil(t, tree)

				return
			}
			assert.NoError(t, err)
			assert.True(t, json.Valid(tree))
		})
	}
}

func TestSymlinkErrorHint(t *testing.T) {
	j := NewJob("Package.swift", testdata.CmdFactoryMock{})
	hint := j.symlinkErrorHint("error: unable to create symlink foo: Permission denied")
	assert.Contains(t, hint, "Developer Mode")
	assert.Contains(t, hint, "Administrator")
	assert.Contains(t, hint, "WSL2")
	assert.Empty(t, j.symlinkErrorHint("error: another failure"))
}
