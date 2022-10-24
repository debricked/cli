package github

import (
	"github.com/debricked/cli/pkg/ci/testdata"
	"github.com/debricked/cli/pkg/ci/util"
	"os"
	"testing"
)

var gitHubActionsEnv = map[string]string{
	"GITHUB_ACTION":     "githubActions",
	"GITHUB_REPOSITORY": "debricked/cli",
	"GITHUB_SHA":        "commit",
	"GITHUB_REF_NAME":   "main",
	"GITHUB_ACTOR":      "viktigpetterr <test@test.com>",
}

func TestIdentify(t *testing.T) {
	ci := Ci{}
	value := os.Getenv(EnvKey)
	if util.EnvKeyIsSet(EnvKey) {
		if !ci.Identify() {
			t.Error("failed to assert that CI was identified")
		}
		_ = os.Unsetenv(EnvKey)
		defer os.Setenv(EnvKey, value)

		if ci.Identify() {
			t.Error("failed to assert that CI was not identified")
		}
	} else {
		if ci.Identify() {
			t.Error("failed to assert that CI was not identified")
		}

		_ = os.Setenv(EnvKey, "value")
		defer os.Unsetenv(EnvKey)

		if !ci.Identify() {
			t.Error("failed to assert that CI identified")
		}
	}
}

func TestParse(t *testing.T) {
	err := testdata.SetUpCiEnv(gitHubActionsEnv)
	if err != nil {
		t.Fatal(err)
	}
	defer testdata.ResetEnv(gitHubActionsEnv)

	ci := Ci{}
	env, _ := ci.Map()
	if env.Filepath != "." {
		t.Error("failed to assert that env contained correct filepath")
	}
	if env.Integration != integration {
		t.Error("failed to assert that env contained correct integration")
	}
	if env.Author != "viktigpetterr <test@test.com>" {
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
