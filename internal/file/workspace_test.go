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
		WorkspacePatterns: []string{
			"package/*",
			"pkg/internal/*",
			"pack/internal/package.json",
			"packages/package_two",
		},
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
		{
			manifestFile: "packages/package_two/package.json",
			expected:     true,
		},
		{
			manifestFile: "packages/package_two/internal/package.json",
			expected:     false,
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

func TestDeeperMatchManifest(t *testing.T) {
	wm := WorkspaceManifest{
		RootManifest: "Src/app/package.json",
		LockFiles:    []string{"Src/app/package-lock.json"},
		WorkspacePatterns: []string{
			"package/*",
			"pkg/internal/*",
			"pack/internal/package.json",
			"packages/package_two",
		},
	}
	cases := []struct {
		manifestFile string
		expected     bool
	}{
		{
			manifestFile: "Src/app/package/test/package.json",
			expected:     true,
		},
		{
			manifestFile: "Src/app/packages/package_two/package.json",
			expected:     true,
		},
		{
			manifestFile: "Src/app/packages/package_two/internal/package.json",
			expected:     false,
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
	workspaces, err := getPackageJSONWorkspaces("testdata/workspace/common/package.json")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(workspaces))
	assert.Equal(t, "packages/*", workspaces[0])
}

func TestGetWorkspacesRare(t *testing.T) {
	workspaces, err := getPackageJSONWorkspaces("testdata/workspace/rare/package.json")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(workspaces))
	assert.Equal(t, "packages/*", workspaces[0])
}

func TestGetWorkspacesNoFile(t *testing.T) {
	_, err := getPackageJSONWorkspaces("testdata/non_existing_folder/package.json")
	assert.Error(t, err)
}
