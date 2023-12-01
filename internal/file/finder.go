package file

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/debricked/cli/internal/client"
	ioFs "github.com/debricked/cli/internal/io"
	"github.com/fatih/color"
)

//go:embed embedded/supported_formats.json
var supportedFormats embed.FS

const SupportedFormatsFallbackFilePath = "embedded/supported_formats.json"
const SupportedFormatsUri = "/api/1.0/open/files/supported-formats"

type IFinder interface {
	GetGroups(rootPath string, exclusions []string, lockfileOnly bool, strictness int) (Groups, error)
	GetSupportedFormats() ([]*CompiledFormat, error)
}

type Finder struct {
	debClient  client.IDebClient
	filesystem ioFs.IFileSystem
}

func NewFinder(c client.IDebClient, fs ioFs.IFileSystem) (*Finder, error) {
	if c == nil {
		return nil, errors.New("client is nil")
	}

	return &Finder{c, fs}, nil
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
	body, err := finder.GetSupportedFormatsJson()
	if err != nil {
		return nil, err
	}

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

func (finder *Finder) GetSupportedFormatsJson() ([]byte, error) {
	finder.debClient.ConfigureClientSettings(false, 3)
	defer finder.debClient.ConfigureClientSettings(true, 15)

	res, err := finder.debClient.Get(SupportedFormatsUri, "application/json")

	if err != nil || res.StatusCode != http.StatusOK {
		fmt.Printf("%s Unable to get supported formats from the server. Using cached data instead.\n", color.YellowString("⚠️"))

		return finder.GetSupportedFormatsFallbackJson()
	}

	defer res.Body.Close()

	return io.ReadAll(res.Body)
}

func (finder *Finder) GetSupportedFormatsFallbackJson() ([]byte, error) {
	jsonFile, err := finder.filesystem.FsOpenEmbed(supportedFormats, SupportedFormatsFallbackFilePath)
	if err != nil {
		return nil, err
	}
	defer finder.filesystem.FsCloseFile(jsonFile)

	jsonData, err := finder.filesystem.FsReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}
