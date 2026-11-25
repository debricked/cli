package mcpclient

import (
	"errors"
	stdfs "io/fs"
	"path/filepath"
	"strings"
	"testing"

	debio "github.com/debricked/cli/internal/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetupInstructionsPreservesAndRefreshesPolicy(t *testing.T) {
	for name, newline := range map[string]string{"LF": "\n", "CRLF": "\r\n"} {
		t.Run(name, func(t *testing.T) {
			fs := debio.FileSystem{}
			path := filepath.Join(t.TempDir(), "instructions.md")
			original := strings.Join([]string{
				"# Existing instructions", "", instructionStart, "outdated policy", instructionEnd, "", "Keep this too.", "",
			}, newline)
			require.NoError(t, fs.FsWriteFile(path, []byte(original), 0644))

			changed, existed, err := SetupInstructions(fs, path, ClientVSCode)

			require.NoError(t, err)
			assert.True(t, changed)
			assert.True(t, existed)
			raw, err := fs.ReadFile(path)
			require.NoError(t, err)
			content := string(raw)
			assert.Equal(t, "# Existing instructions"+newline+newline+instructionStart+newline+
				strings.ReplaceAll(dependencyPolicy, "\n", newline)+newline+instructionEnd+
				newline+newline+"Keep this too."+newline, content)
			_, err = fs.ReadFile(path + BackupSuffix)
			assert.True(t, fs.IsNotExist(err))

			changed, _, err = SetupInstructions(fs, path, ClientVSCode)
			require.NoError(t, err)
			assert.False(t, changed)
			_, err = fs.ReadFile(path + BackupSuffix)
			assert.True(t, fs.IsNotExist(err))
		})
	}
}

func TestSetupInstructionsAppendsToExistingContent(t *testing.T) {
	fs := debio.FileSystem{}
	path := filepath.Join(t.TempDir(), "CLAUDE.md")
	original := "# Project conventions\nKeep existing instructions."
	require.NoError(t, fs.FsWriteFile(path, []byte(original), 0644))

	changed, existed, err := SetupInstructions(fs, path, ClientClaude)

	require.NoError(t, err)
	assert.True(t, changed)
	assert.True(t, existed)
	raw, err := fs.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, original+"\n\n"+instructionStart+"\n"+dependencyPolicy+"\n"+instructionEnd+"\n", string(raw))
}

func TestCursorInstructionsAlwaysApply(t *testing.T) {
	for _, original := range []string{
		"",
		"# Existing rule",
		"---\ndescription: Existing description\nglobs: \"*.go\"\nalwaysApply: false\n---\n\n# Existing rule\n",
	} {
		content, err := upsertInstructions(original, ClientCursor)
		require.NoError(t, err)
		assert.Contains(t, content, "\nalwaysApply: true\n---\n")
		assert.NotContains(t, content, "alwaysApply: false")
		if original != "" {
			assert.Contains(t, content, "# Existing rule")
		}
		if strings.Contains(original, "description:") {
			assert.Contains(t, content, "description: Existing description\nglobs: \"*.go\"")
		}
		updated, err := upsertInstructions(content, ClientCursor)
		require.NoError(t, err)
		assert.Equal(t, content, updated)
	}
}

func TestSetupInstructionsRejectsMalformedManagedContent(t *testing.T) {
	for _, original := range []string{
		instructionStart,
		instructionEnd,
		instructionEnd + "\n" + instructionStart,
		instructionStart + instructionStart + instructionEnd,
		instructionStart + instructionEnd + instructionEnd,
		"---\nalwaysApply: false\n",
	} {
		fs := debio.FileSystem{}
		path := filepath.Join(t.TempDir(), "debricked.mdc")
		require.NoError(t, fs.FsWriteFile(path, []byte(original), 0644))

		_, _, err := SetupInstructions(fs, path, ClientCursor)

		require.Error(t, err)
		raw, err := fs.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, original, string(raw))
		_, err = fs.ReadFile(path + BackupSuffix)
		assert.True(t, fs.IsNotExist(err))
	}
}

type failingInstructionsFS struct {
	debio.IFileSystem
	readErr  error
	mkdirErr error
	writeErr error
}

func (f failingInstructionsFS) ReadFile(path string) ([]byte, error) {
	if f.readErr != nil {
		return nil, f.readErr
	}

	return f.IFileSystem.ReadFile(path)
}

func (f failingInstructionsFS) MkdirAll(path string, perm stdfs.FileMode) error {
	if f.mkdirErr != nil {
		return f.mkdirErr
	}

	return f.IFileSystem.MkdirAll(path, perm)
}

func (f failingInstructionsFS) FsWriteFile(path string, raw []byte, perm stdfs.FileMode) error {
	if f.writeErr != nil {
		return f.writeErr
	}

	return f.IFileSystem.FsWriteFile(path, raw, perm)
}

func TestSetupInstructionsSurfacesFileErrors(t *testing.T) {
	failure := errors.New("permission denied")
	for _, fs := range []failingInstructionsFS{
		{readErr: failure},
		{mkdirErr: failure},
		{writeErr: failure},
	} {
		fs.IFileSystem = debio.FileSystem{}
		path := filepath.Join(t.TempDir(), "instructions.md")
		require.NoError(t, fs.IFileSystem.FsWriteFile(path, []byte("Existing instructions"), 0644))

		changed, _, err := SetupInstructions(fs, path, ClientVSCode)

		assert.ErrorIs(t, err, failure)
		assert.False(t, changed)
		raw, err := fs.IFileSystem.ReadFile(path)
		require.NoError(t, err)
		assert.Equal(t, "Existing instructions", string(raw))
	}
}
