package util

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileContextPlaceholders(t *testing.T) {
	all := &Placeholders{}
	first := NewFileContext(NewProjectScope(all, 1), "a/yarn.lock")
	second := NewFileContext(NewProjectScope(all, 2), "b/paket.lock")

	assert.Equal(t, "debricked-redacted-project-1", first.ProjectPlaceholder("@acme/portal"))
	assert.Equal(t, "debricked-redacted-project-1", first.ProjectPlaceholder("@acme/portal"))
	assert.Equal(t, "debricked-redacted-project-2", first.ProjectPlaceholder("@acme/other"))
	// The same name in another project gets a placeholder of its own.
	assert.Equal(t, "debricked-redacted-project-3", second.ProjectPlaceholder("@acme/portal"))

	assert.Equal(t, "debricked-redacted-file-2-1", second.FilePlaceholder("src/A.cs"))
	assert.Equal(t, "debricked-redacted-file-2-1", second.FilePlaceholder("src/A.cs"))
	assert.Equal(t, "debricked-redacted-file-2-2", second.FilePlaceholder("src/B.cs"))

	// Path placeholders count across projects.
	assert.Equal(t, "file:./redacted-1", first.PathPlaceholder("file:./redacted-%d", "../shared"))
	assert.Equal(t, "file:./redacted-2", second.PathPlaceholder("file:./redacted-%d", "../shared"))

	second.RecordPlaceholder("redacted-1", "packages/worker")

	assert.Equal(t, []MapEntry{
		{Placeholder: "debricked-redacted-project-1", Name: "@acme/portal", Path: "a/yarn.lock"},
		{Placeholder: "debricked-redacted-project-2", Name: "@acme/other", Path: "a/yarn.lock"},
		{Placeholder: "debricked-redacted-project-3", Name: "@acme/portal", Path: "b/paket.lock"},
		{Placeholder: "debricked-redacted-file-2-1", Name: "src/A.cs", Path: "b/paket.lock"},
		{Placeholder: "debricked-redacted-file-2-2", Name: "src/B.cs", Path: "b/paket.lock"},
		{Placeholder: "file:./redacted-1", Name: "../shared", Path: "a/yarn.lock"},
		{Placeholder: "file:./redacted-2", Name: "../shared", Path: "b/paket.lock"},
		{Placeholder: "redacted-1", Name: "packages/worker", Path: "b/paket.lock"},
	}, all.Entries())
}

func TestWriteMap(t *testing.T) {
	path := filepath.Join(t.TempDir(), ".debricked", "sanitize-map.json")
	entries := []MapEntry{{Placeholder: "debricked-redacted-project-1", Name: "@acme/portal", Path: "yarn.lock"}}

	require.NoError(t, WriteMap(path, entries))

	content, err := os.ReadFile(path) // #nosec G304 -- reads the file the test wrote
	require.NoError(t, err)
	var written []MapEntry
	require.NoError(t, json.Unmarshal(content, &written))
	assert.Equal(t, entries, written)
	if runtime.GOOS != "windows" {
		info, err := os.Stat(path)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	}
}

func TestWriteMapWritesAnEmptyListWhenThereAreNoPlaceholders(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sanitize-map.json")

	require.NoError(t, WriteMap(path, nil))

	content, err := os.ReadFile(path) // #nosec G304 -- reads the file the test wrote
	require.NoError(t, err)
	assert.Equal(t, "[]\n", string(content))
}

func TestWriteMapFailsWhenTheFolderCannotBeCreated(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "file")
	require.NoError(t, os.WriteFile(blocker, nil, 0o600))

	assert.Error(t, WriteMap(filepath.Join(blocker, "sanitize-map.json"), nil))
}

func TestIsPackageJSONDependencyList(t *testing.T) {
	assert.True(t, IsPackageJSONDependencyList("devDependencies"))
	assert.False(t, IsPackageJSONDependencyList("scripts"))
}
