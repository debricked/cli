package scan

import (
	"debricked/pkg/client"
	"debricked/pkg/git"
	"strings"
	"testing"
)

func TestNewDebrickedScanner(t *testing.T) {
	var debClient client.IDebClient
	debClient = client.NewDebClient(nil)
	s, err := NewDebrickedScanner(&debClient)

	if err != nil {
		t.Error("failed to assert that no error occurred")
	}

	if s == nil {
		t.Error("failed to assert that scanner was not nil")
	}
}

func TestNewDebrickedScannerWithError(t *testing.T) {
	var debClient client.IDebClient
	s, err := NewDebrickedScanner(&debClient)

	if err == nil {
		t.Error("failed to assert that an error occurred")
	}

	if s != nil {
		t.Error("failed to assert that scanner was nil")
	}

	if !strings.Contains(err.Error(), "failed to initialize the uploader") {
		t.Error("failed to assert error message")
	}
}

func TestScan(t *testing.T) {
	var debClient client.IDebClient
	debClient = client.NewDebClient(nil)
	scanner, _ := NewDebrickedScanner(&debClient)
	directoryPath := "testdata/yarn"
	repositoryName := directoryPath
	commitName := "testdata/yarn-commit"
	opts := DebrickedOptions{
		DirectoryPath:   directoryPath,
		Exclusions:      nil,
		RepositoryName:  repositoryName,
		CommitName:      commitName,
		BranchName:      "",
		CommitAuthor:    "",
		RepositoryUrl:   "",
		IntegrationName: "",
	}
	err := scanner.Scan(opts)
	if err != nil {
		t.Error("failed to assert that scan ran without errors. Error:", err)
	}
}

func TestScanFailingMetaObject(t *testing.T) {
	var debClient client.IDebClient
	debClient = client.NewDebClient(nil)
	scanner, _ := NewDebrickedScanner(&debClient)
	directoryPath := "testdata/yarn"
	opts := DebrickedOptions{
		DirectoryPath:   directoryPath,
		Exclusions:      nil,
		RepositoryName:  "",
		CommitName:      "",
		BranchName:      "",
		CommitAuthor:    "",
		RepositoryUrl:   "",
		IntegrationName: "",
	}
	err := scanner.Scan(opts)
	if err != git.RepositoryNameError {
		t.Error("failed to assert that RepositoryNameError occurred")
	}

	opts.RepositoryName = directoryPath
	err = scanner.Scan(opts)
	if err != git.CommitNameError {
		t.Error("failed to assert that CommitNameError occurred")
	}
}
