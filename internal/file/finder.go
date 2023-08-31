package file

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/debricked/cli/internal/client"
)

type IFinder interface {
	GetGroups(rootPath string, exclusions []string, lockfileOnly bool, strictness int) (Groups, error)
	GetSupportedFormats() ([]*CompiledFormat, error)
}

type Finder struct {
	debClient client.IDebClient
}

func NewFinder(c client.IDebClient) (*Finder, error) {
	if c == nil {
		return nil, errors.New("client is nil")
	}

	return &Finder{c}, nil
}

// GetGroups return all file groups in specified path recursively.
func (finder *Finder) GetGroups(rootPath string, exclusions []string, lockfileOnly bool, strictness int) (Groups, error) {
	var groups Groups

	formats, err := finder.GetSupportedFormats()
	if err != nil {
		return groups, err
	}
	if len(rootPath) == 0 {
		rootPath = filepath.Base("")
	}

	// Traverse files to find dependency file groups
	err = filepath.Walk(
		rootPath,
		func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !fileInfo.IsDir() && !Excluded(exclusions, path) {
				for _, format := range formats {
					if groups.Match(format, path, lockfileOnly) {

						break
					}
				}
			}

			return nil
		},
	)

	groups.FilterGroupsByStrictness(strictness)

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
		} else {
			log.Println(err.Error())
		}
	}

	return compiledDependencyFileFormats, nil
}
