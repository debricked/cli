package github

import (
	"os"
	"testing"

	"github.com/debricked/cli/pkg/ci/testdata"
	"github.com/debricked/cli/pkg/ci/util"
	"github.com/stretchr/testify/assert"
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
		testdata.AssertIdentify(t, ci.Identify, EnvKey)
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
			testdata.SetUpCiEnv(t, c.env)
			defer testdata.ResetEnv(t, c.env)

			env, _ := ci.Map()

			assert.Empty(t, env.Filepath)
			assert.Equal(t, Integration, env.Integration)
			assert.Equal(t, gitHubActionsEnv["GITHUB_ACTOR"], env.Author)
			assert.Equal(t, gitHubActionsEnv["GITHUB_HEAD_REF"], env.Branch)
			assert.Equal(t, "https://github.com/debricked/cli", env.RepositoryUrl)
			assert.Equal(t, gitHubActionsEnv["GITHUB_SHA"], env.Commit)
			assert.Equal(t, "debricked/cli", env.Repository)
		})
	}

}
