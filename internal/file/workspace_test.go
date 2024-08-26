package file

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchManifest(t *testing.T) {
	wm := WorkspaceManifest{
		RootManifest: "package.json",
		LockFiles:    []string{"package-lock.json"},
		Workspaces:   []string{"package/*", "pkg/internal/*", "pack/internal/package.json"},
	}
	cases := []struct {
		manifestFile string
		expected     bool
	}{
		{
			manifestFile: "package/package_one/package.json",
			expected:     true,
		},
		{
			manifestFile: "package_one/package.json",
			expected:     false,
		},
		{
			manifestFile: "pkg/package.json",
			expected:     false,
		},
		{
			manifestFile: "pkg/internal/package.json",
			expected:     true,
		},
	}

	for _, c := range cases {
		t.Run(c.manifestFile, func(t *testing.T) {
			match := (&wm).matchManifest(c.manifestFile)
			fmt.Println(c.manifestFile)
			assert.Equal(t, c.expected, match)
		})
	}
}

func TestGetWorkspaces(t *testing.T) {
	workspaces, err := getWorkspaces("testdata/workspace/package.json")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(workspaces))
	assert.Equal(t, "testdata/workspace/packages/*", workspaces[0])
}

func TestGetWorkspacesNoFile(t *testing.T) {
	_, err := getWorkspaces("testdata/non_existing_folder/package.json")
	assert.Error(t, err)
}
