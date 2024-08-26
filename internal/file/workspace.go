package file

import (
	"encoding/json"
	"os"

	"github.com/becheran/wildmatch-go"
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
		pattern := wildmatch.NewWildMatch(workspacePattern)
		if pattern.IsMatch(manifestPath) {
			return true
		}
	}

	return false
}

func getWorkspaces(rootManifest string) ([]string, error) {
	var npmPackageJson NPMPackageJson

	jsonFile, err := os.Open(rootManifest)
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
