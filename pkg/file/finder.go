package file

import (
	"debricked/pkg/client"
	"encoding/json"
	"errors"
	"github.com/bmatcuk/doublestar/v4"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type Finder struct {
	debClient client.Client
}

func NewFinder(debClient client.Client) (*Finder, error) {
	if debClient == nil {
		return nil, errors.New("DebClient is nil")
	}

	return &Finder{debClient}, nil
}

//GetGroups return all file groups in specified path recursively.
func (finder *Finder) GetGroups(rootPath string, exclusions []string) (Groups, error) {
	var groups Groups

	formats, err := finder.GetSupportedFormats()
	if err != nil {
		return groups, err
	}

	// Traverse files to find dependency file groups
	err = filepath.Walk(
		rootPath,
		func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !fileInfo.IsDir() && !excluded(exclusions, path) {
				for _, format := range formats {
					if groups.Match(format, path) {
						break
					}
				}
			}
			return nil
		},
	)

	return groups, err
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

func excluded(exclusions []string, path string) bool {
	for _, exclusion := range exclusions {
		ex := filepath.Clean(exclusion)
		matched, _ := doublestar.PathMatch(ex, path)
		if matched {
			return true
		}
	}

	return false
}
