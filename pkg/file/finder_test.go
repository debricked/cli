package file

import (
	"debricked/pkg/client"
	"os"
	"strings"
	"testing"
)

var finder *Finder

func setUp(authorized bool) {
	var token string
	if authorized {
		token = os.Getenv("DEBRICKED_TOKEN")
	} else {
		token = "invalid"
	}

	finder, _ = NewFinder(client.NewDebClient(&token))
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
	if len(formats) != 0 {
		t.Error("failed to assert that no formats were found")
	}
	if !strings.Contains(err.Error(), "Unauthorized. Specify access token.") {
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
