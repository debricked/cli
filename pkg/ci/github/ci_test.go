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
	"GITHUB_REF":        "main",
	"GITHUB_ACTOR":      "viktigpetterr <test@test.com>",
	"GITHUB_HEAD_REF":   "main",
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

type parseCase struct {
	name string
	env  map[string]string
}

func TestParse(t *testing.T) {
	var cases []parseCase
	cases = append(cases, parseCase{
		name: "GITHUB_REF with branch",
		env:  gitHubActionsEnv,
	})

	gitHubActionsEnv["GITHUB_REF"] = "refs/tags/main"
	cases = append(cases, parseCase{
		name: "GITHUB_REF with tags",
		env:  gitHubActionsEnv,
	})
	gitHubActionsEnv["GITHUB_REF"] = "refs/heads/main"
	cases = append(cases, parseCase{
		name: "GITHUB_REF with heads",
		env:  gitHubActionsEnv,
	})

	gitHubActionsEnv["GITHUB_REF"] = "refs/pull/18/merge"
	cases = append(cases, parseCase{
		name: "GITHUB_REF with merge",
		env:  gitHubActionsEnv,
	})

	ci := Ci{}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := testdata.SetUpCiEnv(c.env)
			defer testdata.ResetEnv(c.env, t)
			if err != nil {
				t.Fatal(err)
			}

			env, _ := ci.Map()
			if env.Filepath != "." {
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
			if env.RepositoryUrl != "https://github.com/debricked/cli" {
				t.Error("failed to assert that env contained correct repository URL")
			}
			if env.Commit != "commit" {
				t.Error("failed to assert that env contained correct commit")
			}
			if env.Repository != "debricked/cli" {
				t.Error("faield to assert that env contained correct repository")
			}
		})
	}

}
