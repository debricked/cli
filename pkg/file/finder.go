package file

import (
	"debricked/pkg/client"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Finder struct {
	debClient *client.DebClient
}

func NewFinder(debClient *client.DebClient) (*Finder, error) {
	if debClient == nil {
		return nil, errors.New("DebClient is nil")
	}

	return &Finder{debClient}, nil
}

//GetGroups return all file groups in specified path.
func (finder *Finder) GetGroups(path string, ignoredDirs []string) ([]Group, error) {
	formats, err := finder.GetSupportedFormats()
	if err != nil {
		return nil, err
	}
	// Traverse files to find dependency file groups
	var fileGroups []Group
	err = filepath.Walk(
		path,
		func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if ignoredDir(ignoredDirs, path) {
				return filepath.SkipDir
			}

			if !fileInfo.IsDir() {
				for _, format := range formats {
					if format.Match(fileInfo.Name()) {
						fileGroups = append(fileGroups, *NewGroup(path, &Format{}, []string{}))
					}
				}
			}
			return nil
		},
	)

	return fileGroups, err
}

// GetSupportedFormats returns all supported dependency file formats
func (finder *Finder) GetSupportedFormats() ([]*CompiledFormat, error) {
	res, err := finder.debClient.Get("/api/1.0/open/files/supported-formats", "application/json")
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch supported formats")
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	var formats []*Format
	err = json.Unmarshal(body, &formats)
	if err != nil {
		return nil, err
	}

	var compiledDependencyFileFormats []*CompiledFormat
	for _, format := range formats {
		compiledDependencyFileFormat, err := NewCompiledFormat(format)
		if err == nil {
			compiledDependencyFileFormats = append(compiledDependencyFileFormats, compiledDependencyFileFormat)
		}
	}

	return compiledDependencyFileFormats, nil
}

func ignoredDir(ignoredDirs []string, path string) bool {
	for _, exclusion := range ignoredDirs {
		if strings.Contains(path, exclusion) {
			return true
		}
	}

	return false
}
