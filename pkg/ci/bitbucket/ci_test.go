package bitbucket

import (
	"github.com/debricked/cli/pkg/ci/testdata"
	"os"
	"testing"
)

var bitbucketEnv = map[string]string{
	"BITBUCKET_BUILD_NUMBER":    "2",
	"BITBUCKET_REPO_OWNER":      "debricked",
	"BITBUCKET_REPO_SLUG":       "cli",
	"BITBUCKET_COMMIT":          "commit",
	"BITBUCKET_BRANCH":          "main",
	"BITBUCKET_GIT_HTTP_ORIGIN": "https://github.com/debricked/cli",
}

func TestIdentify(t *testing.T) {
	ci := Ci{}

	if ci.Identify() {
		t.Error("failed to assert that CI was not identified")
	}

	_ = os.Setenv(EnvKey, "value")
	defer os.Unsetenv(EnvKey)

	if !ci.Identify() {
		t.Error("failed to assert that CI was identified")
	}
}

func TestParse(t *testing.T) {
	err := testdata.SetUpCiEnv(bitbucketEnv)
	if err != nil {
		t.Fatal(err)
	}
	defer testdata.ResetEnv(bitbucketEnv)

	err = testdata.SetUpGitRepository()
	if err != nil {
		t.Fatal("failed to initialize repository", err)
	}
	defer testdata.TearDownGitRepository()

	ci := Ci{}
	env, err := ci.Map()
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}
	if env.Filepath != "." {
		t.Error("failed to assert that env contained correct filepath")
	}
	if env.Integration != Integration {
		t.Error("failed to assert that env contained correct integration")
	}
	if len(env.Author) == 0 {
		t.Error("failed to assert that env contained correct author")
	}
	if env.Branch != "main" {
		t.Error("failed to assert that env contained correct branch")
	}
	if env.RepositoryUrl != "https://github.com/debricked/cli" {
		t.Error("failed to assert that env contained correct repository URL")
	}
	if env.Commit != "commit" {
		t.Error("failed to assert that env contained correct commit")
	}
	if env.Repository != "debricked/cli" {
		t.Error("faield to assert that env contained correct repository")
	}
}
