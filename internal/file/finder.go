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

type DebrickedOptions struct {
	RootPath     string
	Exclusions   []string
	Inclusions   []string
	LockFileOnly bool
	Strictness   int
}

type IFinder interface {
	GetGroups(options DebrickedOptions) (Groups, error)
	GetSupportedFormats() ([]*CompiledFormat, error)
	GetConfigPath(rootPath string, exclusions []string, inclusions []string) string
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

func (finder *Finder) GetConfigPath(rootPath string, exclusions []string, inclusions []string) string {
	var configPath string

	if len(rootPath) == 0 {
		rootPath = filepath.Base("")
	}
	err := filepath.Walk(
		rootPath,
		func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !fileInfo.IsDir() && !Excluded(exclusions, inclusions, path) {
				if filepath.Base(path) == "debricked-config.yaml" {
					configPath = path
				}
			}

			return nil
		},
	)
	if err != nil {
		return ""
	}

	return configPath
}

func isCompressed(filename string) bool {
	compressionExtensions := map[string]struct{}{
		".gz":  {},
		".zip": {},
		".tar": {},
		".rar": {},
		".bz2": {},
		".xz":  {},
		".7z":  {},
	}
	myExt := filepath.Ext(filename)
	_, compressed := compressionExtensions[myExt]

	return compressed
}

func (finder *Finder) GetIncludedGroups(formats []*CompiledFormat, options DebrickedOptions) (Groups, error) {
	// NOTE: inefficient because it walks into excluded directories
	var groups Groups
	err := filepath.Walk(
		options.RootPath,
		func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			var excluded = Excluded(options.Exclusions, options.Inclusions, path)

			if !fileInfo.IsDir() && !excluded {
				for _, format := range formats {
					if groups.Match(format, path, options.LockFileOnly) {

						break
					}
				}
			}

			return nil
		},
	)

	return groups, err
}

func (finder *Finder) GetExcludedGroups(formats []*CompiledFormat, options DebrickedOptions) (Groups, []string, error) {
	var excludedGroups Groups
	var excludedFiles []string
	err := filepath.Walk(
		options.RootPath,
		func(path string, fileInfo os.FileInfo, err error) error {
			if err != nil {

				return err
			}
			if isCompressed(path) {
				excludedFiles = append(excludedFiles, path)
			} else if !fileInfo.IsDir() {
				for _, format := range formats {
					if excludedGroups.Match(format, path, options.LockFileOnly) {
						excludedFiles = append(excludedFiles, path)

						break
					}
				}
			}

			return nil
		},
	)

	return excludedGroups, excludedFiles, err
}

func reportExclusions(excludedFiles []string) {
	if len(excludedFiles) > 0 {
		containsCompressedFile := false
		fmt.Println("The following files were excluded, resulting in no dependency files found.")
		for _, file := range excludedFiles {
			if !containsCompressedFile && isCompressed(file) {
				containsCompressedFile = true
			}
			fmt.Println(file)
		}
		if containsCompressedFile {
			fmt.Println("Compressed file found, but contained files cannot be scanned. Decompress to scan content.")
		}
	} else {
		fmt.Println("No dependency file matches found with current configuration.")
	}
	fmt.Println("Change the inclusion and exclusion options if a file or directory was missed.")

}

// GetGroups return all file groups in specified path recursively.
func (finder *Finder) GetGroups(options DebrickedOptions) (Groups, error) {
	var groups Groups

	formats, err := finder.GetSupportedFormats()
	if err != nil {

		return groups, err
	}
	if len(options.RootPath) == 0 {
		options.RootPath = filepath.Base("")
	}

	// Traverse files to find dependency file groups
	groups, err = finder.GetIncludedGroups(formats, options)
	if len(groups.groups) == 0 {
		// No dependencies found. (should rarely happen)
		// Traverse again to see if dependency or zip files were excluded.
		_, excludedFiles, excludedErr := finder.GetExcludedGroups(formats, options)
		reportExclusions(excludedFiles)
		if excludedErr != nil {

			return groups, err
		}
	}
	groups.AddWorkspaceLockFiles()
	groups.FilterGroupsByStrictness(options.Strictness)

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

	sbtEntry := &Format{
		ManifestFileRegex: "^build\\.sbt$",
		DocumentationUrl:  "https://docs.debricked.com/overview/language-support/scala-sbt",
		LockFileRegexes:   []string{""},
	}

	formats = append(formats, sbtEntry)

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
