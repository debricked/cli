package file

import (
	"bytes"
	"debricked/pkg/client"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

type debClientMock struct{}

func (mock *debClientMock) Post(_ string, _ string, _ *bytes.Buffer) (*http.Response, error) {
	return nil, nil
}

var authorized bool

func (mock *debClientMock) Get(_ string, _ string) (*http.Response, error) {
	var statusCode int
	var body io.ReadCloser = nil
	if authorized {
		statusCode = http.StatusOK
		formatsBytes, _ := json.Marshal(formatsMock)
		body = ioutil.NopCloser(strings.NewReader(string(formatsBytes)))
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
	if err == nil {
		t.Error("failed to assert that error occurred")
	}
	if finder != nil {
		t.Error("failed to assert that finder was nil")
	}

	if !strings.Contains(err.Error(), "DebClient is nil") {
		t.Error("failed to assert error message")
	}

	finder, err = NewFinder(client.NewDebClient(nil))
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}
	if finder == nil {
		t.Error("failed to assert that finder was not nil")
	}
}

func TestGetSupportedFormats(t *testing.T) {
	setUp(true)
	formats, err := finder.GetSupportedFormats()
	if err != nil {
		t.Fatal("failed to assert that no error occurred. Error:", err)
	}
	if len(formats) == 0 {
		t.Error("failed to assert that there is formats")
	}
	for _, format := range formats {
		hasContent := format.Regex != nil || len(format.LockFileRegexes) > 0
		if !hasContent {
			t.Error("failed to assert that format had content")
		}
	}
}

func TestGetSupportedFormatsFailed(t *testing.T) {
	setUp(false)
	formats, err := finder.GetSupportedFormats()
	if len(formats) > 0 {
		t.Error("failed to assert that no formats were found")
	}
	if !strings.Contains(err.Error(), "failed to fetch supported formats") {
		t.Error("failed to assert error message")
	}
}

func TestGetGroups(t *testing.T) {
	setUp(true)
	directoryPath := "."
	fileGroups, err := finder.GetGroups(directoryPath, []string{"testdata/go"})
	if err != nil {
		t.Fatal("failed to assert that no error occurred. Error:", err)
	}
	if len(fileGroups) == 0 {
		t.Error("failed to assert that there is at least one format")
	}
	for _, fileGroup := range fileGroups {
		hasContent := fileGroup.CompiledFormat != nil && (strings.Contains(fileGroup.FilePath, directoryPath) || len(fileGroup.RelatedFiles) > 0)
		if !hasContent {
			t.Error("failed to assert that format had content")
		}
	}
}

func TestIgnoredDir(t *testing.T) {
	dir := "composer"
	ignoredDirs := []string{dir}
	files, _ := os.ReadDir("testdata")
	for _, file := range files {
		if file.Name() == dir {
			if !ignoredDir(ignoredDirs, file.Name()) {
				t.Error("failed to assert that directory was not ignored")
			}
		} else if file.IsDir() {
			if ignoredDir(ignoredDirs, file.Name()) {
				t.Error("failed to assert that directory was not ignored")
			}
		}
	}
}
