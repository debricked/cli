package circleci

import (
	"testing"

	"github.com/debricked/cli/pkg/ci/env"
	"github.com/debricked/cli/pkg/ci/testdata"
	"github.com/stretchr/testify/assert"
)

const (
	debrickedUrl = "https://github.com/debricked/cli"
)

var circleCiEnv = map[string]string{
	"CIRCLECI":                "circleci",
	"CIRCLE_PROJECT_USERNAME": "debricked",
	"CIRCLE_PROJECT_REPONAME": "cli",
	"CIRCLE_SHA1":             "commit",
	"CIRCLE_BRANCH":           "main",
	"CIRCLE_REPOSITORY_URL":   "https://github.com/debricked/cli.git",
}

func TestIdentify(t *testing.T) {
	testdata.AssertIdentify(t, Ci{}.Identify, EnvKey)
}

func TestParse(t *testing.T) {
	testdata.SetUpCiEnv(t, circleCiEnv)
	defer testdata.ResetEnv(t, circleCiEnv)

	cwd := testdata.SetUpGitRepository(t, true)
	defer testdata.TearDownGitRepository(cwd, t)

	ci := Ci{}
	e, err := ci.Map()
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}

	assertEnv(e, t)
}

func TestMapRepositoryUrl(t *testing.T) {
	ci := Ci{}
	cases := map[string]string{
		"https://github.com/debricked/cli.git":    debrickedUrl,
		"http://gitlab.com/debricked/cli.git":     "http://gitlab.com/debricked/cli",
		"http://gitlab.com/debricked/sub/cli.git": "http://gitlab.com/debricked/sub/cli",
		"git@github.com:debricked/cli.git":        debrickedUrl,
		"git@gitlab.com:debricked/cli.git":        "https://gitlab.com/debricked/cli",
		"tcp@scm.com:debricked/sub/cli.git":       "tcp@scm.com:debricked/sub/cli.git",
	}
	for gitUrl, assertion := range cases {
		t.Run(gitUrl, func(t *testing.T) {
			repository := ci.MapRepositoryUrl(gitUrl)
			assert.Equal(t, assertion, repository)
		})
	}
}

func assertEnv(env env.Env, t *testing.T) {
	assert.Empty(t, env.Filepath)
	assert.Equal(t, Integration, env.Integration)
	assert.NotEmpty(t, env.Author)
	assert.Equal(t, circleCiEnv["CIRCLE_BRANCH"], env.Branch)
	assert.Equal(t, debrickedUrl, env.RepositoryUrl)
	assert.Equal(t, circleCiEnv["CIRCLE_SHA1"], env.Commit)
	assert.Equal(t, "debricked/cli", env.Repository)
}
