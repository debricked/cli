package scan

import (
	"github.com/debricked/cli/pkg/ci"
	"github.com/debricked/cli/pkg/ci/argo"
	"github.com/debricked/cli/pkg/ci/azure"
	"github.com/debricked/cli/pkg/ci/bitbucket"
	"github.com/debricked/cli/pkg/ci/buildkite"
	"github.com/debricked/cli/pkg/ci/circleci"
	"github.com/debricked/cli/pkg/ci/env"
	"github.com/debricked/cli/pkg/ci/github"
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
	path := "testdata/yarn"
	repositoryName := path
	commitName := "testdata/yarn-commit"
	opts := DebrickedOptions{
		Path:            path,
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
	path := "testdata/yarn"
	opts := DebrickedOptions{
		Path:            path,
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

	opts.RepositoryName = path
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
	path := "."
	opts := DebrickedOptions{
		Path:            path,
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

func TestMapEnvToOptions(t *testing.T) {
	dOptionsTemplate := DebrickedOptions{
		Path:            "path",
		Exclusions:      nil,
		RepositoryName:  "repository",
		CommitName:      "commit",
		BranchName:      "branch",
		CommitAuthor:    "author",
		RepositoryUrl:   "url",
		IntegrationName: "CLI",
	}

	cases := []struct {
		name     string
		template DebrickedOptions
		opts     DebrickedOptions
		env      env.Env
	}{
		{
			name:     "No env",
			template: dOptionsTemplate,
			opts:     dOptionsTemplate,
			env: env.Env{
				Repository:    "",
				Commit:        "",
				Branch:        "",
				Author:        "",
				RepositoryUrl: "",
				Integration:   "",
				Filepath:      "",
			},
		},
		{
			name: "CI env set",
			template: DebrickedOptions{
				Path:            "env-path",
				Exclusions:      nil,
				RepositoryName:  "env-repository",
				CommitName:      "env-commit",
				BranchName:      "env-branch",
				CommitAuthor:    "author",
				RepositoryUrl:   "env-url",
				IntegrationName: github.Integration,
			},
			opts: DebrickedOptions{
				Path:            "input-path",
				Exclusions:      nil,
				RepositoryName:  "",
				CommitName:      "",
				BranchName:      "",
				CommitAuthor:    "author",
				RepositoryUrl:   "",
				IntegrationName: "CLI",
			},
			env: env.Env{
				Repository:    "env-repository",
				Commit:        "env-commit",
				Branch:        "env-branch",
				Author:        "env-author",
				RepositoryUrl: "env-url",
				Integration:   github.Integration,
				Filepath:      "env-path",
			},
		},
		{
			name: "CI env set without directory path",
			template: DebrickedOptions{
				Path:            "input-path",
				Exclusions:      nil,
				RepositoryName:  "env-repository",
				CommitName:      "env-commit",
				BranchName:      "env-branch",
				CommitAuthor:    "author",
				RepositoryUrl:   "env-url",
				IntegrationName: github.Integration,
			},
			opts: DebrickedOptions{
				Path:            "input-path",
				Exclusions:      nil,
				RepositoryName:  "",
				CommitName:      "",
				BranchName:      "",
				CommitAuthor:    "author",
				RepositoryUrl:   "",
				IntegrationName: "CLI",
			},
			env: env.Env{
				Repository:    "env-repository",
				Commit:        "env-commit",
				Branch:        "env-branch",
				Author:        "env-author",
				RepositoryUrl: "env-url",
				Integration:   github.Integration,
				Filepath:      "",
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			MapEnvToOptions(&c.opts, c.env)
			strings.EqualFold(c.opts.Path, c.template.Path)
			if !strings.EqualFold(c.opts.Path, c.template.Path) {
				t.Errorf("Failed to assert that %s was equal to %s", c.opts.Path, c.template.Path)
			}
			if c.opts.Exclusions != nil {
				t.Errorf("Failed to assert that Exclusions was nil")
			}
			if c.opts.RepositoryName != c.template.RepositoryName {
				t.Errorf("Failed to assert that %s was equal to %s", c.opts.RepositoryName, c.template.RepositoryName)
			}
			if c.opts.CommitName != c.template.CommitName {
				t.Errorf("Failed to assert that %s was equal to %s", c.opts.CommitName, c.template.CommitName)
			}
			if c.opts.BranchName != c.template.BranchName {
				t.Errorf("Failed to assert that %s was equal to %s", c.opts.BranchName, c.template.BranchName)
			}
			if c.opts.CommitAuthor != c.template.CommitAuthor {
				t.Errorf("Failed to assert that %s was equal to %s", c.opts.CommitAuthor, c.template.CommitAuthor)
			}
			if c.opts.RepositoryUrl != c.template.RepositoryUrl {
				t.Errorf("Failed to assert that %s was equal to %s", c.opts.RepositoryUrl, c.template.RepositoryUrl)
			}
			if c.opts.IntegrationName != c.template.IntegrationName {
				t.Errorf("Failed to assert that %s was equal to %s", c.opts.IntegrationName, c.template.IntegrationName)
			}
		})
	}
}
