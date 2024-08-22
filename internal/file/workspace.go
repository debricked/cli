package file

import (
	"encoding/json"
	"github.com/bmatcuk/doublestar/v4"
	"os"
)

type WorkspaceManifest struct {
	RootManifest string
	LockFile     string
	Workspaces   []string
}

type NPMPackageJson struct {
	Name       string   `json:"name"`
	Workspaces []string `json:"workspaces"`
}

func (workspaceManifest *WorkspaceManifest) matchManifest(manifestPath string) bool {
	for _, workspacePattern := range workspaceManifest.Workspaces {
		matched, _ := doublestar.Match(workspacePattern, manifestPath)
		if matched {

			return true
		}
	}

	return false
}

func getWorkspaces(packageJson string) ([]string, error) {
	var npmPackageJson NPMPackageJson

	jsonFile, err := os.Open(packageJson)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&npmPackageJson)
	if err != nil {
		return nil, err
	}

	return npmPackageJson.Workspaces, nil
}
