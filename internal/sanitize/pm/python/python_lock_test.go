package python

import (
	"strings"
	"testing"

	"github.com/debricked/cli/internal/sanitize/pm/util"
	"github.com/pelletier/go-toml/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSanitizePipfileLock(t *testing.T) {
	before := `{
  "_meta": {"sources": [{"url": "https://user:secret@example.test/simple"}]},
  "default": {
    "first": {"version": "==1.2.3", "git": "https://token@example.test/repo", "ref": "main"},
    "second": {"git": "https://token@example.test/repo", "ref": "main"},
    "empty": {"markers": "python_version > '3.8'"}
  },
  "develop": {"dev-package": {"version": "==4.5.6", "index": "internal"}}
}`

	sanitized, err := sanitizePipfileLock([]byte(before), nil)
	require.NoError(t, err)
	assert.Equal(t, `{
    "default": {
        "first": {
            "version": "==1.2.3",
            "git": "redacted"
        },
        "second": {},
        "empty": {}
    },
    "develop": {
        "dev-package": {
            "version": "==4.5.6"
        }
    }
}
`, string(sanitized))
	assert.Less(t, strings.Index(string(sanitized), `"first"`), strings.Index(string(sanitized), `"second"`))
}

func TestSanitizePoetryLock(t *testing.T) {
	before := `[[package]]
name = "first"
version = "1.0"
description = "private details"

[package.dependencies]
zlib = "^1.0"
alpha = { version = "2.0", source = "internal" }

[[package]]
name = "second"
version = "2.0"
source = { url = "https://user:secret@example.test" }

[metadata]
lock-version = "2.0"
`

	sanitized, err := sanitizePoetryLock([]byte(before), nil)
	require.NoError(t, err)
	assert.NotContains(t, string(sanitized), "metadata")
	assert.NotContains(t, string(sanitized), "example.test")
	var lock map[string]any
	require.NoError(t, toml.Unmarshal(sanitized, &lock))
	packages, ok := lock["package"].([]any)
	require.True(t, ok)
	require.Len(t, packages, 2)
	first, ok := packages[0].(map[string]any)
	require.True(t, ok)
	second, ok := packages[1].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "first", first["name"])
	assert.Equal(t, "second", second["name"])
	assert.Equal(t, map[string]any{"zlib": "*", "alpha": "*"}, first["dependencies"])
}

func TestSanitizeUVLock(t *testing.T) {
	before := `version = 1

[[package]]
name = "my-project"
version = "1.2.3"
source = { virtual = "." }

[[package]]
name = "first"
version = "4.5.6"
source = { registry = "https://user:secret@example.test/simple" }
sdist = { url = "https://example.test/first.tar.gz", hash = "secret" }
wheels = [{ url = "https://example.test/first.whl", hash = "secret" }]
dependencies = [{ name = "my-project", version = "1.2.3" }, { name = "second", marker = "python_version > '3.8'" }]

[package.metadata]
requires-dist = ["private"]

[manifest]
members = ["."]
`
	scope := util.NewProjectScope(&util.Placeholders{}, 1)
	context := util.NewFileContext(scope, "uv.lock")

	sanitized, err := sanitizeUVLock([]byte(before), context)
	require.NoError(t, err)
	assert.NotContains(t, string(sanitized), "source")
	assert.NotContains(t, string(sanitized), "sdist")
	assert.NotContains(t, string(sanitized), "wheels")
	assert.NotContains(t, string(sanitized), "metadata")
	assert.NotContains(t, string(sanitized), "manifest")
	assert.NotContains(t, string(sanitized), "example.test")

	var lock map[string]any
	require.NoError(t, toml.Unmarshal(sanitized, &lock))
	packages, ok := lock["package"].([]any)
	require.True(t, ok)
	project, ok := packages[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "debricked-redacted-project-1", project["name"])
	assert.Equal(t, "0.0.0", project["version"])
	first, ok := packages[1].(map[string]any)
	require.True(t, ok)
	dependencies, ok := first["dependencies"].([]any)
	require.True(t, ok)
	projectDependency, ok := dependencies[0].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "debricked-redacted-project-1", projectDependency["name"])
	secondDependency, ok := dependencies[1].(map[string]any)
	require.True(t, ok)
	assert.Equal(t, "second", secondDependency["name"])
}

func TestPythonPmOnlyHandlesLockFiles(t *testing.T) {
	pm := NewPm()
	assert.Equal(t, Name, pm.Name())
	lockFiles := pm.LockFiles()
	assert.Len(t, lockFiles, 3)
	assert.Contains(t, lockFiles, "Pipfile.lock")
	assert.Contains(t, lockFiles, "poetry.lock")
	assert.Contains(t, lockFiles, "uv.lock")
	assert.Empty(t, pm.Manifests())
}

func TestSanitizePythonLockRejectsInvalidInput(t *testing.T) {
	for name, sanitizer := range map[string]util.FileSanitizer{
		"Pipfile.lock": sanitizePipfileLock,
		"poetry.lock":  sanitizePoetryLock,
		"uv.lock":      sanitizeUVLock,
	} {
		t.Run(name, func(t *testing.T) {
			_, err := sanitizer([]byte("not a lockfile"), nil)
			assert.ErrorIs(t, err, util.ErrInvalidFormat)
		})
	}
}
