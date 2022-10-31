package travis

import (
	"github.com/debricked/cli/pkg/ci/testdata"
	"os"
	"testing"
)

var travisEnv = map[string]string{
	"TRAVIS_REPO_SLUG": "debricked/cli",
	"TRAVIS_BRANCH":    "main",
	"TRAVIS_COMMIT":    "commit",
	"TRAVIS_BUILD_DIR": ".",
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
	err := testdata.SetUpCiEnv(travisEnv)
	defer testdata.ResetEnv(travisEnv)
	if err != nil {
		t.Error(err)
	}

	err = testdata.SetUpGitRepository()
	defer testdata.TearDownGitRepository()
	if err != nil {
		t.Error("failed to initialize repository", err)
	}

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
