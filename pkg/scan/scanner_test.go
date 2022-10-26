package scan

import (
	"github.com/debricked/cli/pkg/ci"
	"github.com/debricked/cli/pkg/ci/argo"
	"github.com/debricked/cli/pkg/ci/azure"
	"github.com/debricked/cli/pkg/ci/bitbucket"
	"github.com/debricked/cli/pkg/ci/buildkite"
	"github.com/debricked/cli/pkg/ci/circleci"
	"github.com/debricked/cli/pkg/ci/gitlab"
	"github.com/debricked/cli/pkg/ci/travis"
	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/git"
	"github.com/debricked/cli/pkg/upload"
	"strings"
	"testing"
)

func TestNewDebrickedScanner(t *testing.T) {
	var debClient client.IDebClient
	debClient = client.NewDebClient(nil)
	var ciService ci.IService
	s, err := NewDebrickedScanner(&debClient, ciService)

	if err != nil {
		t.Error("failed to assert that no error occurred")
	}

	if s == nil {
		t.Error("failed to assert that scanner was not nil")
	}
}

func TestNewDebrickedScannerWithError(t *testing.T) {
	var debClient client.IDebClient
	var ciService ci.IService
	s, err := NewDebrickedScanner(&debClient, ciService)

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
	var ciService ci.IService
	ciService = ci.NewService(nil)
	scanner, _ := NewDebrickedScanner(&debClient, ciService)
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
	var ciService ci.IService
	ciService = ci.NewService([]ci.ICi{
		argo.Ci{},
		azure.Ci{},
		bitbucket.Ci{},
		buildkite.Ci{},
		circleci.Ci{},
		//github.Ci{}, Since GitHub actions is used, this ICi is ignored
		gitlab.Ci{},
		travis.Ci{},
	})
	scanner, _ := NewDebrickedScanner(&debClient, ciService)
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

func TestScanFailingNoFiles(t *testing.T) {
	var debClient client.IDebClient
	debClient = client.NewDebClient(nil)
	var ciService ci.IService
	ciService = ci.NewService([]ci.ICi{
		argo.Ci{},
		azure.Ci{},
		bitbucket.Ci{},
		buildkite.Ci{},
		circleci.Ci{},
		//github.Ci{}, Since GitHub actions is used, this ICi is ignored
		gitlab.Ci{},
		travis.Ci{},
	})
	scanner, _ := NewDebrickedScanner(&debClient, ciService)
	directoryPath := "."
	opts := DebrickedOptions{
		DirectoryPath:   directoryPath,
		Exclusions:      []string{"testdata/**"},
		RepositoryName:  "name",
		CommitName:      "commit",
		BranchName:      "branch",
		CommitAuthor:    "",
		RepositoryUrl:   "",
		IntegrationName: "",
	}
	err := scanner.Scan(opts)
	if err != upload.NoFilesErr {
		t.Error("failed to assert that error NoFilesErr occurred")
	}
}

func TestScanBadOpts(t *testing.T) {
	var c client.IDebClient
	scanner, _ := NewDebrickedScanner(&c, nil)
	var opts IOptions
	err := scanner.Scan(opts)
	if err != BadOptsErr {
		t.Error("failed to assert that BadOptsErr occurred")
	}
}
