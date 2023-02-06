package scan

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"github.com/debricked/cli/pkg/client/testdata"
	"github.com/debricked/cli/pkg/file"
	"github.com/debricked/cli/pkg/git"
	"github.com/debricked/cli/pkg/upload"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const (
	testdataYarn = "testdata/yarn"
)

func TestNewDebrickedScanner(t *testing.T) {
	var debClient client.IDebClient = client.NewDebClient(nil)
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
	clientMock := testdata.NewDebClientMock()
	addMockedFormatsResponse(clientMock)
	addMockedFileUploadResponse(clientMock)
	addMockedFinishResponse(clientMock, http.StatusNoContent)
	addMockedStatusResponse(clientMock, http.StatusOK, 50)
	addMockedStatusResponse(clientMock, http.StatusOK, 100)
	debClient = clientMock

	var ciService ci.IService = ci.NewService(nil)

	scanner, _ := NewDebrickedScanner(&debClient, ciService)

	path := testdataYarn
	repositoryName := path
	commitName := "testdata/yarn-commit"
	cwd, _ := os.Getwd()
	// reset working directory that has been manipulated in scanner.Scan
	defer resetWd(t, cwd)
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

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := scanner.Scan(opts)

	_ = w.Close()
	output, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	if err != nil {
		t.Error("failed to assert that scan ran without errors. Error:", err)
	}

	outputAssertions := []string{
		"Working directory: /",
		"cli/pkg/scan/testdata/yarn\n",
		"Successfully uploaded:  yarn.lock\n",
		"Successfully concluded upload\n",
		"Scanning...",
		"0% |",
		"50% |",
		"100% |",
		"32m✔",
		"0 vulnerabilities found\n\n",
		"For full details, visit:",
	}
	for _, assertion := range outputAssertions {
		if !strings.Contains(string(output), assertion) {
			t.Errorf("failed to assert %s in output", assertion)
		}
	}
}

func TestScanFailingMetaObject(t *testing.T) {
	var debClient client.IDebClient = client.NewDebClient(nil)
	var ciService ci.IService = ci.NewService([]ci.ICi{
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
	cwd, _ := os.Getwd()
	path := testdataYarn
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
	// reset working directory that has been manipulated in scanner.Scan
	_ = os.Chdir(cwd)

	opts.RepositoryName = path
	err = scanner.Scan(opts)
	if err != git.CommitNameError {
		t.Error("failed to assert that CommitNameError occurred")
	}
	// reset working directory that has been manipulated in scanner.Scan
	_ = os.Chdir(cwd)
}

func TestScanFailingNoFiles(t *testing.T) {
	var debClient client.IDebClient
	clientMock := testdata.NewDebClientMock()
	addMockedFormatsResponse(clientMock)
	debClient = clientMock
	var ciService ci.IService = ci.NewService([]ci.ICi{
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
	opts := DebrickedOptions{
		Path:            "",
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

func TestScanEmptyResult(t *testing.T) {
	var debClient client.IDebClient
	clientMock := testdata.NewDebClientMock()
	addMockedFormatsResponse(clientMock)
	addMockedFileUploadResponse(clientMock)
	addMockedFinishResponse(clientMock, http.StatusNoContent)
	addMockedStatusResponse(clientMock, http.StatusOK, 50)
	// Create mocked scan result response, 201 is returned when the queue time are too long
	addMockedStatusResponse(clientMock, http.StatusCreated, 0)

	debClient = clientMock

	var ciService ci.IService = ci.NewService(nil)
	scanner, _ := NewDebrickedScanner(&debClient, ciService)
	path := testdataYarn
	repositoryName := path
	commitName := "testdata/yarn-commit"
	cwd, _ := os.Getwd()
	// reset working directory that has been manipulated in scanner.Scan
	defer resetWd(t, cwd)

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

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := scanner.Scan(opts)

	_ = w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	existsMessageInCMDOutput := strings.Contains(
		string(out),
		"Progress polling terminated due to long scan times. Please try again later")

	if err != nil || !existsMessageInCMDOutput {
		t.Error("failed to assert that scan ran without errors. Error:", err)
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

func TestSetWorkingDirectory(t *testing.T) {
	absPath, _ := filepath.Abs("")
	cases := []struct {
		name        string
		opts        DebrickedOptions
		errMessages []string
	}{
		{
			name: "empty path",
			opts: DebrickedOptions{Path: ""},
		},
		{
			name: "absolute path",
			opts: DebrickedOptions{Path: absPath},
		},
		{
			name: "relative path",
			opts: DebrickedOptions{Path: ".."},
		},
		{
			name: "current working directory",
			opts: DebrickedOptions{Path: "."},
		},
		{
			name:        "bad path",
			opts:        DebrickedOptions{Path: "bad-path"},
			errMessages: []string{"no such file or directory", "The system cannot find the file specified"},
		},
	}
	cwd, _ := os.Getwd()
	defer resetWd(t, cwd)
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := SetWorkingDirectory(&c.opts)

			if len(c.errMessages) > 0 {
				containsCorrectErrMsg := false
				for _, errMsg := range c.errMessages {
					containsCorrectErrMsg = containsCorrectErrMsg || strings.Contains(err.Error(), errMsg)
				}
				if !containsCorrectErrMsg {
					t.Errorf("failed to assert that error message contained either of: %s or %s. Got: %s", c.errMessages[0], c.errMessages[1], err.Error())
				}
			} else {
				if len(c.opts.Path) != 0 {
					t.Errorf("failed to assert that Path was empty. Got: %s", c.opts.Path)
				}
			}
		})
	}
}

func addMockedFormatsResponse(clientMock *testdata.DebClientMock) {
	formats := []file.Format{{
		Regex:           "",
		LockFileRegexes: []string{"yarn\\.lock"},
	}}
	formatsBytes, _ := json.Marshal(formats)
	formatsMockRes := testdata.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: io.NopCloser(bytes.NewReader(formatsBytes)),
	}
	clientMock.AddMockUriResponse("/api/1.0/open/files/supported-formats", formatsMockRes)
}

func addMockedFileUploadResponse(clientMock *testdata.DebClientMock) {
	uploadMockRes := testdata.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: io.NopCloser(strings.NewReader("{\"ciUploadId\": 1}")),
	}
	clientMock.AddMockUriResponse("/api/1.0/open/uploads/dependencies/files", uploadMockRes)
}

func addMockedFinishResponse(clientMock *testdata.DebClientMock, statusCode int) {
	finishMockRes := testdata.MockResponse{
		StatusCode:   statusCode,
		ResponseBody: io.NopCloser(strings.NewReader("{}")),
	}
	clientMock.AddMockUriResponse("/api/1.0/open/finishes/dependencies/files/uploads", finishMockRes)
}

func addMockedStatusResponse(clientMock *testdata.DebClientMock, statusCode int, progress int) {
	finishMockRes := testdata.MockResponse{
		StatusCode:   statusCode,
		ResponseBody: io.NopCloser(strings.NewReader(fmt.Sprintf(`{"progress": %d}`, progress))),
	}
	clientMock.AddMockUriResponse("/api/1.0/open/ci/upload/status", finishMockRes)
}

func resetWd(t *testing.T, wd string) {
	err := os.Chdir(wd)
	if err != nil {
		t.Fatal("Can not read the directory: ", wd)
	}
}
