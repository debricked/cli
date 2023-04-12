package file

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/debricked/cli/internal/client/testdata"
	"github.com/stretchr/testify/assert"
)

type debClientMock struct{}

func (mock *debClientMock) Post(_ string, _ string, _ *bytes.Buffer) (*http.Response, error) {
	return &http.Response{}, nil
}

var authorized bool

func (mock *debClientMock) Get(_ string, _ string) (*http.Response, error) {
	var statusCode int
	var body io.ReadCloser = nil
	if authorized {
		statusCode = http.StatusOK
		formatsBytes, _ := json.Marshal(formatsMock)
		body = io.NopCloser(strings.NewReader(string(formatsBytes)))
	} else {
		statusCode = http.StatusForbidden
	}
	res := http.Response{
		Status:           "",
		StatusCode:       statusCode,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             body,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}

	return &res, nil
}

var finder *Finder

func setUp(auth bool) {
	finder, _ = NewFinder(&debClientMock{})
	authorized = auth
}

func TestNewFinder(t *testing.T) {
	finder, err := NewFinder(nil)

	assert.NotNil(t, err)
	assert.Nil(t, finder)
	assert.ErrorContains(t, err, "client is nil")

	finder, err = NewFinder(testdata.NewDebClientMock())
	assert.Nil(t, err)
	assert.NotNil(t, finder)
}

func TestGetSupportedFormats(t *testing.T) {
	setUp(true)
	formats, err := finder.GetSupportedFormats()
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(formats), 1)
	for _, format := range formats {
		hasContent := format.Regex != nil || len(format.LockFileRegexes) > 0
		assert.True(t, hasContent, "failed to assert that format had content")
	}
}

func TestGetSupportedFormatsFailed(t *testing.T) {
	setUp(false)
	formats, err := finder.GetSupportedFormats()
	assert.Error(t, err)
	assert.ErrorContains(t, err, "failed to fetch supported formats")
	assert.Empty(t, formats)
}

func TestGetGroups(t *testing.T) {
	setUp(true)
	path := ""

	exclusions := []string{"testdata/go/*.mod", "testdata/misc/**"}
	excludedFiles := []string{"testdata/go/go.mod", "testdata/misc/requirements.txt"}
	const nbrOfGroups = 4

	fileGroups, err := finder.GetGroups(path, exclusions, false, StrictAll)

	assert.NoError(t, err)
	assert.Equalf(t, nbrOfGroups, fileGroups.Size(), "failed to assert that %d groups were created. %d was found", nbrOfGroups, fileGroups.Size())

	for _, fileGroup := range fileGroups.ToSlice() {
		hasContent := fileGroup.CompiledFormat != nil && (strings.Contains(fileGroup.FilePath, path) || len(fileGroup.RelatedFiles) > 0)
		assert.True(t, hasContent, "failed to assert that format had content")

		groupFiles := fileGroup.RelatedFiles
		groupFiles = append(groupFiles, fileGroup.FilePath)
		for _, groupFile := range groupFiles {
			for _, exFile := range excludedFiles {
				assert.NotEqualf(t, groupFile, exFile, "failed to assert that file was excluded")
			}
		}
	}
}

func TestExclude(t *testing.T) {
	var files []string
	_ = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			files = append(files, path)

			return nil
		})

	cases := []struct {
		name               string
		exclusions         []string
		expectedExclusions []string
	}{
		{
			name:               "NoExclusions",
			exclusions:         []string{},
			expectedExclusions: []string{},
		},
		{
			name:               "InvalidFileExclusion",
			exclusions:         []string{"composer.json"},
			expectedExclusions: []string{},
		},
		{
			name:               "FileExclusionWithDoublestar",
			exclusions:         []string{"**/composer.json"},
			expectedExclusions: []string{"composer.json"},
		},
		{
			name:               "DirectoryExclusion",
			exclusions:         []string{"*/composer/*"},
			expectedExclusions: []string{"composer.json", "composer.lock"},
		},
		{
			name:               "DirectoryExclusionWithRelPath",
			exclusions:         []string{"testdata/go/*"},
			expectedExclusions: []string{"go.mod"},
		},
		{
			name:               "ExtensionExclusionWithWildcardAndDoublestar",
			exclusions:         []string{"**/*.mod"},
			expectedExclusions: []string{"go.mod"},
		},
		{
			name:               "DirectoryExclusionWithDoublestar",
			exclusions:         []string{"**/yarn/**"},
			expectedExclusions: []string{"yarn", "yarn.lock"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var excludedFiles []string
			for _, file := range files {
				if excluded(c.exclusions, file) {
					excludedFiles = append(excludedFiles, file)
				}
			}

			assert.Equal(t, len(c.expectedExclusions), len(excludedFiles), "failed to assert that the same number of files were ignored")

			for _, file := range excludedFiles {
				baseName := filepath.Base(file)
				asserted := false
				for _, expectedExcludedFile := range c.expectedExclusions {
					if baseName == expectedExcludedFile {
						asserted = true

						break
					}
				}

				assert.Truef(t, asserted, "%s ignored when it should pass", file)
			}
		})
	}
}

func TestGetGroupsWithOnlyLockFiles(t *testing.T) {
	setUp(true)
	path := "testdata/misc"
	const nbrOfGroups = 1
	fileGroups, err := finder.GetGroups(path, []string{"**/requirements*.txt"}, false, StrictAll)
	assert.NoError(t, err)
	assert.Equalf(t, nbrOfGroups, fileGroups.Size(), "failed to assert that %d groups were created. %d was found", nbrOfGroups, fileGroups.Size())

	fileGroup := fileGroups.groups[0]
	assert.False(t, fileGroup.HasFile(), "failed to assert that file group lacked file")
	assert.Len(t, fileGroup.RelatedFiles, 1, "failed to assert that there was one related file")

	file := fileGroup.GetAllFiles()[0]

	assert.Contains(t, file, "Cargo.lock", "failed to assert that the related file was Cargo.lock")
}

func TestGetGroupsWithTwoFileMatchesInSameDir(t *testing.T) {
	setUp(true)
	path := "testdata/pip"
	const nbrOfGroups = 2
	fileGroups, err := finder.GetGroups(path, []string{}, false, StrictAll)
	assert.NoError(t, err)
	assert.Equalf(t, nbrOfGroups, fileGroups.Size(), "failed to assert that %d groups were created. %d was found", nbrOfGroups, fileGroups.Size())

	var files []string
	for _, fg := range fileGroups.groups {
		assert.Len(t, fg.RelatedFiles, 1)
		files = append(files, filepath.Base(fg.FilePath))
		relatedFile := fg.RelatedFiles[0]
		if strings.Contains(fg.FilePath, "requirements.txt") {
			assert.Contains(t, relatedFile, "requirements.txt.pip.debricked.lock")
		} else {
			assert.Contains(t, relatedFile, "requirements-dev.txt.pip.debricked.lock")
		}
	}

	assert.Contains(t, files, "requirements.txt")
	assert.Contains(t, files, "requirements-dev.txt")
}

func TestGetGroupsWithStrictFlag(t *testing.T) {
	setUp(true)
	cases := []struct {
		name                   string
		strictness             int
		testedGroupIndex       int
		expectedNumberOfGroups int
		expectedManifestFile   string
		expectedLockFiles      []string
	}{
		{
			name:                   "StrictnessSetTo0",
			strictness:             StrictAll,
			testedGroupIndex:       3,
			expectedNumberOfGroups: 7,
			expectedManifestFile:   "requirements.txt",
			expectedLockFiles:      []string{},
		},
		{
			name:                   "StrictnessSetTo1",
			strictness:             StrictLockAndPairs,
			testedGroupIndex:       1,
			expectedNumberOfGroups: 5,
			expectedManifestFile:   "",
			expectedLockFiles: []string{
				"Cargo.lock",
			},
		},
		{
			name:                   "StrictnessSetTo2",
			strictness:             StrictPairs,
			testedGroupIndex:       0,
			expectedNumberOfGroups: 3,
			expectedManifestFile:   "composer.json",
			expectedLockFiles: []string{
				"composer.lock",
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			path := "testdata"
			fileGroups, err := finder.GetGroups(path, []string{}, false, c.strictness)
			fileGroup := fileGroups.groups[c.testedGroupIndex]

			assert.Nilf(t, err, "failed to assert that no error occurred. Error: %s", err)
			assert.NotNilf(t, fileGroup, "failed to find group with index: %d", c.testedGroupIndex)
			assert.Equalf(
				t,
				c.expectedNumberOfGroups,
				fileGroups.Size(),
				"failed to assert that %d groups were created. %d were found",
				c.expectedNumberOfGroups,
				fileGroups.Size(),
			)
			assert.Containsf(
				t,
				fileGroup.FilePath,
				c.expectedManifestFile,
				"actual manifest file %s doesn't match expected %s",
				fileGroup.FilePath,
				c.expectedManifestFile,
			)

			if len(c.expectedLockFiles) > 0 {
				for i := range c.expectedLockFiles {
					assert.Containsf(
						t,
						fileGroup.RelatedFiles[i],
						c.expectedLockFiles[i],
						"actual lock file %s doesn't match expected %s",
						fileGroup.RelatedFiles[i],
						c.expectedLockFiles[i],
					)
				}
			}
		})
	}
}
