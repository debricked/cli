package file

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/becheran/wildmatch-go"
)

type WorkspaceManifest struct {
	RootManifest      string
	LockFiles         []string
	WorkspacePatterns []string
}

// NPM Workspaces docs: https://docs.npmjs.com/cli/v10/configuring-npm/package-json#workspaces
// Yarn Workspaces docs: https://yarnpkg.com/features/workspaces
type PackageJSON struct {
	Name       string   `json:"name"`
	Workspaces []string `json:"workspaces"`
}

type NestledPackageJSON struct {
	Workspaces struct {
		Packages []string `json:"packages"`
	} `json:"workspaces"`
} // Rare format

func (workspaceManifest *WorkspaceManifest) matchManifest(manifestPath string) bool {
	manifestPath = filepath.ToSlash(manifestPath) // Normalize
	relativeManifestPath, err := filepath.Rel(
		filepath.ToSlash(filepath.Dir(workspaceManifest.RootManifest)),
		manifestPath,
	)
	if err != nil {
		relativeManifestPath = manifestPath
	}
	relativeManifestPath = filepath.ToSlash(relativeManifestPath)
	for _, workspacePattern := range workspaceManifest.WorkspacePatterns {
		workspacePattern = filepath.ToSlash(workspacePattern)
		pattern := wildmatch.NewWildMatch(workspacePattern)
		if pattern.IsMatch(relativeManifestPath) {
			return true
		}
		if workspacePattern == filepath.ToSlash(filepath.Dir(relativeManifestPath)) {
			return true
		} // Check if specific directory match
	}

	return false
}

func getPackageJSONWorkspaces(rootManifest string) ([]string, error) {
	var packageJson PackageJSON

	jsonData, err := os.ReadFile(rootManifest)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &packageJson)
	if err != nil { // Try rare format instead
		var nestledPackageJSON NestledPackageJSON
		err = json.Unmarshal(jsonData, &nestledPackageJSON)
		if err == nil {
			return nestledPackageJSON.Workspaces.Packages, nil
		}

		return nil, err
	}

	return packageJson.Workspaces, nil
}
