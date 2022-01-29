package scan

import (
	"debricked/pkg/client"
	"debricked/pkg/git"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestNewScanCmd(t *testing.T) {
	cmd := NewScanCmd(client.NewDebClient(nil))
	flags := cmd.Flags()
	flagAssertions := map[string]string{
		"repository":     "r",
		"commit":         "c",
		"branch":         "b",
		"author":         "a",
		"repository-url": "u",
		"integration":    "i",
	}
	for name, shorthand := range flagAssertions {
		flag := flags.Lookup(name)
		if flag == nil {
			t.Error(fmt.Sprintf("failed to assert that %s flag was set", name))
		}
		if flag.Shorthand != shorthand {
			t.Error(fmt.Sprintf("failed to assert that %s flag shorthand %s was set correctly", name, shorthand))
		}
	}
}

func TestValidateArgs(t *testing.T) {
	err := validateArgs(nil, []string{"/"})
	if err != nil {
		t.Error("failed to assert that no error occurred. Error: ", err.Error())
	}
}

func TestValidateArgsMissingArg(t *testing.T) {
	err := validateArgs(nil, []string{})
	if err == nil {
		t.Error("failed to assert that an error occurred")
	}
	if !strings.Contains(err.Error(), "requires directory path") {
		t.Error("failed assert error")
	}
}

func TestValidateArgsInvalidArg(t *testing.T) {
	err := validateArgs(nil, []string{"invalid-path"})
	if err == nil {
		t.Error("failed to assert that an error occurred")
	}
	if !strings.Contains(err.Error(), "invalid directory path specified") {
		t.Error("failed assert error")
	}
}

func TestRun(t *testing.T) {
	debClient = client.NewDebClient(nil)
	repositoryName = "testdata/yarn"
	commitName = "testdata/yarn-commit"
	err := run(nil, []string{"testdata/yarn"})
	if err != nil {
		t.Fatal("failed to assert that no error occurred. Error: ", err.Error())
	}
}

func TestRunMissingRepositoryProperties(t *testing.T) {
	debClient = client.NewDebClient(nil)
	repositoryName = ""
	err := run(nil, []string{"testdata/yarn"})
	if err == nil {
		t.Fatal("failed to assert that an error occurred")
	}
	if !strings.Contains(err.Error(), "failed to find repository name. Please use --repository flag") {
		t.Error("failed to assert error message")
	}
}

func TestGetSupportedFormats(t *testing.T) {
	debClient = client.NewDebClient(nil)
	formats, err := getSupportedFormats()
	if err != nil {
		t.Fatal("failed to assert that no error occurred. Error: ", err.Error())
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

func TestFindFileGroups(t *testing.T) {
	debClient = client.NewDebClient(nil)
	directoryPath := "../../.."
	fileGroups, err := findFileGroups(directoryPath)
	if err != nil {
		t.Fatal("failed to assert that no error occurred. Error: ", err.Error())
	}
	if len(fileGroups) == 0 {
		t.Error("failed to assert that there is formats")
	}
	for _, fileGroup := range fileGroups {
		hasContent := fileGroup.Format != nil && (strings.Contains(fileGroup.FilePath, directoryPath) || len(fileGroup.RelatedFiles) > 0)
		if !hasContent {
			t.Error("failed to assert that format had content")
		}
	}
}

func TestIgnoredDir(t *testing.T) {
	dir := "testdata"
	ignoredDirs = []string{dir}
	files, _ := os.ReadDir("../..")
	for _, file := range files {
		fileInfo, _ := file.Info()
		if file.Name() == dir {
			if !ignoredDir(fileInfo) {
				t.Error("failed to assert that directory was not ignored")
			}
		} else if file.IsDir() {
			if ignoredDir(fileInfo) {
				t.Error("failed to assert that directory was not ignored")
			}
		}
	}
}

func TestScan(t *testing.T) {
	directoryPath := "testdata/yarn"
	repositoryName = directoryPath
	commitName = "testdata/yarn-commit"
	metaObject, err := git.NewMetaObject(directoryPath, repositoryName, commitName, "", "", "")
	err = scan(directoryPath, metaObject, []string{})
	if err != nil {
		t.Error("failed to assert that scan ran without errors. Error: ", err.Error())
	}
}
