package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchManifest(t *testing.T) {
	wm := WorkspaceManifest{
		RootManifest: "package.json",
		LockFile:     "package-lock.json",
		Workspaces:   []string{"package/*"},
	}
	match := (&wm).matchManifest("package/package_one/package.json")
	assert.True(t, match)
}

func TestGetWorkspaces(t *testing.T) {
	workspaces, err := getWorkspaces("testdata/workspace/package.json")
	assert.NoError(t, err)
	assert.Equal(t, len(workspaces), 1)
	assert.Equal(t, workspaces[0], "packages/*")
}

func TestGetWorkspacesNoFile(t *testing.T) {
	_, err := getWorkspaces("testdata/non_existing_folder/package.json")
	assert.Error(t, err)
}
