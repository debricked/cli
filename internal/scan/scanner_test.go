package scan

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/debricked/cli/internal/callgraph"
	callgraphTestdata "github.com/debricked/cli/internal/callgraph/testdata"
	"github.com/debricked/cli/internal/ci"
	"github.com/debricked/cli/internal/ci/argo"
	"github.com/debricked/cli/internal/ci/azure"
	"github.com/debricked/cli/internal/ci/bitbucket"
	"github.com/debricked/cli/internal/ci/buildkite"
	"github.com/debricked/cli/internal/ci/circleci"
	"github.com/debricked/cli/internal/ci/env"
	"github.com/debricked/cli/internal/ci/github"
	"github.com/debricked/cli/internal/ci/gitlab"
	"github.com/debricked/cli/internal/ci/travis"
	"github.com/debricked/cli/internal/client"
	"github.com/debricked/cli/internal/client/testdata"
	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/fingerprint"
	"github.com/debricked/cli/internal/git"
	ioFs "github.com/debricked/cli/internal/io"
	"github.com/debricked/cli/internal/resolution"
	resolveTestdata "github.com/debricked/cli/internal/resolution/testdata"
	"github.com/debricked/cli/internal/upload"
	"github.com/stretchr/testify/assert"
)

const windowsOS = "windows"

var testdataNpm = filepath.Join("testdata", "npm")

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

func TestNewDebrickedScanner(t *testing.T) {
	var debClient client.IDebClient
	var cis ci.IService
	var finder file.IFinder
	var uploader upload.IUploader
	var resolver resolution.IResolver
	var fingerprint fingerprint.IFingerprint
	var generator callgraph.IGenerator
	s := NewDebrickedScanner(&debClient, finder, uploader, cis, resolver, fingerprint, generator)

	assert.NotNil(t, s)
}

func TestScan(t *testing.T) {
	if runtime.GOOS == windowsOS {
		t.Skipf("TestScan is skipped due to Windows env")
	}
	clientMock := testdata.NewDebClientMock()
	addMockedFormatsResponse(clientMock, "package\\.json")
	addMockedFileUploadResponse(clientMock)
	addMockedFinishResponse(clientMock, http.StatusNoContent)
	addMockedStatusResponse(clientMock, http.StatusOK, 50)
	addMockedStatusResponse(clientMock, http.StatusOK, 100)
	scanner := makeScanner(clientMock, nil, nil)

	path := testdataNpm
	repositoryName := path
	cwd, _ := os.Getwd()
	// reset working directory that has been manipulated in scanner.Scan
	defer resetWd(t, cwd)
	opts := DebrickedOptions{
		Path:                     path,
		Exclusions:               nil,
		RepositoryName:           repositoryName,
		CommitName:               "commit",
		Fingerprint:              false,
		BranchName:               "",
		CommitAuthor:             "",
		RepositoryUrl:            "",
		IntegrationName:          "",
		CallGraphUploadTimeout:   10 * 60,
		CallGraphGenerateTimeout: 10 * 60,
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
		"Successfully uploaded",
		"package.json\n",
		"Successfully initialized scan\n",
		"Scanning...",
		"0% |",
		"50% |",
		"100% |",
		"32m✔",
		"0 vulnerabilities found\n\n",
		"For full details, visit:",
	}
	for _, assertion := range outputAssertions {
		assert.Contains(t, string(output), assertion)
	}
}

func TestScanWithJsonPath(t *testing.T) {
	if runtime.GOOS == windowsOS {
		t.Skipf("TestScan is skipped due to Windows env")
	}
	clientMock := testdata.NewDebClientMock()
	addMockedFormatsResponse(clientMock, "package\\.json")
	addMockedFileUploadResponse(clientMock)
	addMockedFinishResponse(clientMock, http.StatusNoContent)
	addMockedStatusResponse(clientMock, http.StatusOK, 50)
	addMockedStatusResponse(clientMock, http.StatusOK, 100)
	scanner := makeScanner(clientMock, nil, nil)

	path := testdataNpm
	repositoryName := path
	cwd, _ := os.Getwd()

	// reset working directory that has been manipulated in scanner.Scan
	defer resetWd(t, cwd)
	opts := DebrickedOptions{
		Path:                     path,
		Exclusions:               nil,
		RepositoryName:           repositoryName,
		CommitName:               "commit",
		BranchName:               "",
		CommitAuthor:             "",
		RepositoryUrl:            "",
		IntegrationName:          "",
		CallGraphUploadTimeout:   10 * 60,
		CallGraphGenerateTimeout: 10 * 60,
		JsonFilePath:             "result.json",
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
		"Successfully uploaded",
		"package.json",
		"Successfully initialized scan\n",
		"Scanning...",
		"0% |",
		"50% |",
		"100% |",
		"32m✔",
		"0 vulnerabilities found\n\n",
		"For full details, visit:",
	}
	for _, assertion := range outputAssertions {
		assert.Contains(t, string(output), assertion)
	}
	assert.FileExists(t, filepath.Join(cwd, path, "result.json"))
}

func TestScanFailingMetaObject(t *testing.T) {
	var debClient client.IDebClient = testdata.NewDebClientMock()
	scanner := NewDebrickedScanner(&debClient, nil, nil, ciService, nil, nil, nil)
	cwd, _ := os.Getwd()
	path := testdataNpm
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

	assert.ErrorIs(t, err, git.RepositoryNameError)
	_ = os.Chdir(cwd)

	opts.RepositoryName = path
	err = scanner.Scan(opts)

	assert.ErrorIs(t, err, git.CommitNameError)
	_ = os.Chdir(cwd)
}

func TestScanFailingNoFiles(t *testing.T) {
	clientMock := testdata.NewDebClientMock()
	addMockedFormatsResponse(clientMock, "package\\.json")
	scanner := makeScanner(clientMock, nil, nil)
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

	assert.ErrorIs(t, err, upload.NoFilesErr)
}

func TestScanBadOpts(t *testing.T) {
	var c client.IDebClient
	scanner := NewDebrickedScanner(&c, nil, nil, nil, nil, nil, nil)
	var opts IOptions

	err := scanner.Scan(opts)

	assert.ErrorIs(t, err, BadOptsErr)
}

func TestScanEmptyResult(t *testing.T) {
	if runtime.GOOS == windowsOS {
		t.Skipf("TestScan is skipped due to Windows env")
	}
	clientMock := testdata.NewDebClientMock()
	addMockedFormatsResponse(clientMock, "package\\.json")
	addMockedFileUploadResponse(clientMock)
	addMockedFinishResponse(clientMock, http.StatusNoContent)
	addMockedStatusResponse(clientMock, http.StatusOK, 50)
	// Create mocked scan result response, 201 is returned when the queue time are too long
	addMockedStatusResponse(clientMock, http.StatusCreated, 0)

	scanner := makeScanner(clientMock, nil, nil)
	path := testdataNpm
	repositoryName := path
	commitName := testdataNpm
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

	assert.NoError(t, err, "failed to assert that scan ran without errors")
	assert.True(t, existsMessageInCMDOutput, "failed to assert that scan ran without errors")
}

func TestScanInCiWithPathSet(t *testing.T) {
	var debClient client.IDebClient = testdata.NewDebClientMock()
	scanner := NewDebrickedScanner(&debClient, nil, nil, ciService, nil, nil, nil)
	cwd, _ := os.Getwd()
	defer resetWd(t, cwd)
	path := testdataNpm
	_ = os.Setenv("GITLAB_CI", "gitlab")
	_ = os.Setenv("CI_PROJECT_DIR", ".")
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
	assert.ErrorIs(t, git.RepositoryNameError, err)
	cwd, _ = os.Getwd()
	assert.Contains(t, cwd, testdataNpm)
}

func TestScanWithResolve(t *testing.T) {
	clientMock := testdata.NewDebClientMock()
	addMockedFormatsResponse(clientMock, "yarn\\.lock")
	addMockedFileUploadResponse(clientMock)
	addMockedFinishResponse(clientMock, http.StatusNoContent)
	addMockedStatusResponse(clientMock, http.StatusOK, 100)

	resolverMock := resolveTestdata.ResolverMock{}
	resolverMock.SetFiles([]string{"yarn.lock"})

	scanner := makeScanner(clientMock, &resolverMock, nil)

	cwd, _ := os.Getwd()
	defer resetWd(t, cwd)
	// Clean up resolution must be done before wd reset, otherwise files cannot be deleted
	defer cleanUpResolution(t, resolverMock)

	path := testdataNpm
	repositoryName := path
	commitName := "testdata/npm-commit"
	opts := DebrickedOptions{
		Path:            path,
		Resolve:         true,
		Exclusions:      nil,
		RepositoryName:  repositoryName,
		CommitName:      commitName,
		BranchName:      "",
		CommitAuthor:    "",
		RepositoryUrl:   "",
		IntegrationName: "",
	}
	err := scanner.Scan(opts)
	assert.NoError(t, err)
	cwd, _ = os.Getwd()
	assert.Contains(t, cwd, path)
}

func TestScanWithResolveErr(t *testing.T) {
	clientMock := testdata.NewDebClientMock()
	resolutionErr := errors.New("resolution-error")
	scanner := makeScanner(clientMock, &resolveTestdata.ResolverMock{Err: resolutionErr}, nil)
	cwd, _ := os.Getwd()
	defer resetWd(t, cwd)

	path := testdataNpm
	repositoryName := path
	commitName := "testdata/yarn-commit"
	opts := DebrickedOptions{
		Path:           path,
		Resolve:        true,
		Exclusions:     nil,
		RepositoryName: repositoryName,
		CommitName:     commitName,
	}
	err := scanner.Scan(opts)
	assert.ErrorIs(t, err, resolutionErr)
}

// TestScanWithResolveErr tests that the scan is not aborted if the resolution fails
var dOptionsTemplate = DebrickedOptions{
	Path:            "path",
	Exclusions:      nil,
	Inclusions:      nil,
	RepositoryName:  "repository",
	CommitName:      "commit",
	BranchName:      "branch",
	CommitAuthor:    "author",
	RepositoryUrl:   "url",
	IntegrationName: "CLI",
}

var cases = []struct {
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
			Path:            "",
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
	{
		name: "CI env set with directory path",
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
			Filepath:      "env-path",
		},
	},
}

func TestMapEnvToOptions(t *testing.T) {

	for _, co := range cases {
		c := co
		t.Run(c.name, func(t *testing.T) {
			MapEnvToOptions(&c.opts, c.env)
			assert.Equal(t, c.template.Path, c.opts.Path)
			assert.Nil(t, c.opts.Exclusions)
			assert.Nil(t, c.opts.Inclusions)
			assert.Equal(t, c.template.RepositoryName, c.opts.RepositoryName)
			assert.Equal(t, c.template.CommitName, c.opts.CommitName)
			assert.Equal(t, c.template.BranchName, c.opts.BranchName)
			assert.Equal(t, c.template.CommitAuthor, c.opts.CommitAuthor)
			assert.Equal(t, c.template.RepositoryUrl, c.opts.RepositoryUrl)
			assert.Equal(t, c.template.IntegrationName, c.opts.IntegrationName)
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
	for _, co := range cases {
		c := co
		t.Run(c.name, func(t *testing.T) {
			err := SetWorkingDirectory(&c.opts)
			defer resetWd(t, cwd)

			if len(c.errMessages) > 0 {
				containsCorrectErrMsg := false
				for _, errMsg := range c.errMessages {
					containsCorrectErrMsg = containsCorrectErrMsg || strings.Contains(err.Error(), errMsg)
				}
				assert.Truef(t, containsCorrectErrMsg, "failed to assert that error message contained either of: %s or %s. Got: %s", c.errMessages[0], c.errMessages[1], err.Error())
			} else {
				assert.Lenf(t, c.opts.Path, 0, "failed to assert that Path was empty. Got: %s", c.opts.Path)
			}
		})
	}
}

func TestScanServiceDowntime(t *testing.T) {
	clientMock := testdata.NewDebClientMock()
	clientMock.SetServiceUp(false)

	scanner := makeScanner(clientMock, nil, nil)
	path := testdataNpm
	repositoryName := path
	commitName := "testdata/yarn-commit"
	cwd, _ := os.Getwd()
	// reset working directory that has been manipulated in scanner.Scan
	defer resetWd(t, cwd)
	opts := DebrickedOptions{
		Path:           path,
		Exclusions:     nil,
		RepositoryName: repositoryName,
		CommitName:     commitName,
		Fingerprint:    false,
		PassOnTimeOut:  true,
	}

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := scanner.Scan(opts)

	_ = w.Close()
	output, _ := io.ReadAll(r)
	os.Stdout = rescueStdout

	assert.ErrorContains(t, err, client.NoResErr.Error())
	assert.Contains(t, string(output), client.SupportedFormatsFallbackError.Error())
	resetWd(t, cwd)

	opts.PassOnTimeOut = false
	rescueStdout = os.Stdout
	os.Stdout = w

	err = scanner.Scan(opts)

	_ = w.Close()
	os.Stdout = rescueStdout
	assert.ErrorContains(t, err, client.NoResErr.Error())
}

func TestScanWithFingerprint(t *testing.T) {
	if runtime.GOOS == windowsOS {
		t.Skipf("TestScan is skipped due to Windows env")
	}
	clientMock := testdata.NewDebClientMock()
	clientMock.SetEnterpriseCustomer(true)
	addMockedFormatsResponse(clientMock, "yarn\\.lock")
	addMockedFileUploadResponse(clientMock)
	addMockedFinishResponse(clientMock, http.StatusNoContent)
	addMockedStatusResponse(clientMock, http.StatusOK, 100)

	resolverMock := resolveTestdata.ResolverMock{}
	resolverMock.SetFiles([]string{"yarn.lock"})

	scanner := makeScanner(clientMock, &resolverMock, nil)
	scanner.fingerprint = fingerprint.NewFingerprinter()

	cwd, _ := os.Getwd()
	defer resetWd(t, cwd)
	// Clean up resolution must be done before wd reset, otherwise files cannot be deleted
	defer cleanUpResolution(t, resolverMock)

	path := testdataNpm
	repositoryName := path
	commitName := "testdata/npm-commit-fingerprint"
	opts := DebrickedOptions{
		Path:            path,
		Resolve:         true,
		Fingerprint:     true,
		Exclusions:      nil,
		Inclusions:      nil,
		RepositoryName:  repositoryName,
		CommitName:      commitName,
		BranchName:      "",
		CommitAuthor:    "",
		RepositoryUrl:   "",
		IntegrationName: "",
	}
	err := scanner.Scan(opts)
	assert.NoError(t, err)
	cwd, _ = os.Getwd()
	assert.Contains(t, cwd, path)
}

func TestScanWithFingerprintNoEnterprise(t *testing.T) {
	if runtime.GOOS == windowsOS {
		t.Skipf("TestScan is skipped due to Windows env")
	}
	clientMock := testdata.NewDebClientMock()
	addMockedFormatsResponse(clientMock, "yarn\\.lock")
	addMockedFileUploadResponse(clientMock)
	addMockedFinishResponse(clientMock, http.StatusNoContent)
	addMockedStatusResponse(clientMock, http.StatusOK, 100)

	resolverMock := resolveTestdata.ResolverMock{}
	resolverMock.SetFiles([]string{"yarn.lock"})

	scanner := makeScanner(clientMock, &resolverMock, nil)
	scanner.fingerprint = fingerprint.NewFingerprinter()

	cwd, _ := os.Getwd()
	defer resetWd(t, cwd)
	defer cleanUpResolution(t, resolverMock)

	path := testdataNpm
	repositoryName := path
	commitName := "testdata/npm-commit-fingerprint"
	opts := DebrickedOptions{
		Path:            path,
		Resolve:         true,
		Fingerprint:     true,
		Exclusions:      nil,
		Inclusions:      nil,
		RepositoryName:  repositoryName,
		CommitName:      commitName,
		BranchName:      "",
		CommitAuthor:    "",
		RepositoryUrl:   "",
		IntegrationName: "",
	}
	err := scanner.Scan(opts)
	assert.NoError(t, err)
	cwd, _ = os.Getwd()
	assert.Contains(t, cwd, path)
}

func TestScanWithCallgraph(t *testing.T) {
	if runtime.GOOS == windowsOS {
		t.Skipf("TestScan is skipped due to Windows env")
	}
	clientMock := testdata.NewDebClientMock()
	addMockedFormatsResponse(clientMock, "yarn\\.lock")
	addMockedFileUploadResponse(clientMock)
	addMockedFinishResponse(clientMock, http.StatusNoContent)
	addMockedStatusResponse(clientMock, http.StatusOK, 100)

	generatorMock := callgraphTestdata.GeneratorMock{}

	scanner := makeScanner(clientMock, nil, &generatorMock)
	scanner.fingerprint = fingerprint.NewFingerprinter()

	cwd, _ := os.Getwd()
	defer resetWd(t, cwd)

	path := testdataNpm
	repositoryName := path
	commitName := "testdata/npm-commit-callgraph"
	opts := DebrickedOptions{
		Path:            path,
		Resolve:         false,
		Fingerprint:     false,
		CallGraph:       true,
		Exclusions:      nil,
		Inclusions:      nil,
		RepositoryName:  repositoryName,
		CommitName:      commitName,
		BranchName:      "",
		CommitAuthor:    "",
		RepositoryUrl:   "",
		IntegrationName: "",
	}
	err := scanner.Scan(opts)
	assert.ErrorContains(t, err, "failed to find dependency files")
	cwd, _ = os.Getwd()
	assert.Contains(t, cwd, path)
}

func addMockedFormatsResponse(clientMock *testdata.DebClientMock, regex string) {
	formats := []file.Format{{
		ManifestFileRegex: "",
		LockFileRegexes:   []string{regex},
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

func makeScanner(clientMock *testdata.DebClientMock, resolverMock *resolveTestdata.ResolverMock, generatorMock *callgraphTestdata.GeneratorMock) *DebrickedScanner {
	var debClient client.IDebClient = clientMock

	var finder file.IFinder
	finder, _ = file.NewFinder(debClient, ioFs.FileSystem{})

	var uploader upload.IUploader
	uploader, _ = upload.NewUploader(debClient)

	var cis ci.IService = ci.NewService(nil)

	return NewDebrickedScanner(&debClient, finder, uploader, cis, resolverMock, nil, generatorMock)
}

func cleanUpResolution(t *testing.T, resolverMock resolveTestdata.ResolverMock) {
	err := resolverMock.CleanUp()
	if err != nil {
		t.Error(err)
	}
}
