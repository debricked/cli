package gitlab

import (
	"github.com/debricked/cli/pkg/ci/testdata"
	"os"
	"testing"
)

var gitLabEnv = map[string]string{
	"GITLAB_CI":          "gitlab",
	"CI_PROJECT_PATH":    "debricked/cli",
	"CI_COMMIT_SHA":      "commit",
	"CI_COMMIT_REF_NAME": "main",
	"CI_PROJECT_DIR":     "/",
	"CI_PROJECT_URL":     "https://gitlab.com/debricked/cli",
	"CI_COMMIT_AUTHOR":   "viktigpetterr <test@test.com>",
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
	err := testdata.SetUpCiEnv(gitLabEnv)
	if err != nil {
		t.Fatal(err)
	}
	defer testdata.ResetEnv(gitLabEnv)

	ci := Ci{}
	env, _ := ci.Map()
	if env.Filepath != "/" {
		t.Error("failed to assert that env contained correct filepath")
	}
	if env.Integration != Integration {
		t.Error("failed to assert that env contained correct integration")
	}
	if env.Author != "viktigpetterr <test@test.com>" {
		t.Error("failed to assert that env contained correct author")
	}
	if env.Branch != "main" {
		t.Error("failed to assert that env contained correct branch")
	}
	if env.RepositoryUrl != "https://gitlab.com/debricked/cli" {
		t.Error("failed to assert that env contained correct repository URL")
	}
	if env.Commit != "commit" {
		t.Error("failed to assert that env contained correct commit")
	}
	if env.Repository != "debricked/cli" {
		t.Error("faield to assert that env contained correct repository")
	}
}
