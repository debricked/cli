package argo

import (
	"testing"

	"github.com/debricked/cli/internal/ci/env"
	"github.com/debricked/cli/internal/ci/testdata"
	"github.com/stretchr/testify/assert"
)

const (
	debrickedCli = "debricked/cli"
	debrickedUrl = "https://github.com/debricked/cli"
)

var argoEnv = map[string]string{
	"DEBRICKED_GIT_URL":       "https://github.com/debricked/cli.git",
	"BASE_DIRECTORY":          "/",
	"ARGO_AGENT_TASK_WORKERS": "16",
}

func TestIdentify(t *testing.T) {
	testdata.AssertIdentify(t, Ci{}.Identify, EnvKey)
}

func TestParse(t *testing.T) {
	testdata.SetUpCiEnv(t, argoEnv)
	defer testdata.ResetEnv(t, argoEnv)

	cwd := testdata.SetUpGitRepository(t, true)
	defer testdata.TearDownGitRepository(cwd, t)

	ci := Ci{}
	e, err := ci.Map()
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}
	assertEnv(e, t)
}

func TestMapRepository(t *testing.T) {
	ci := Ci{}
	cases := []string{
		"https://github.com/debricked/cli.git",
		"http://gitlab.com/debricked/cli.git",
		"http://scm.com/debricked/cli.git",
		"git@github.com:debricked/cli.git",
		"git@gitlab.com:debricked/cli.git",
		"tcp@scm.com:debricked/cli.git",
	}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			repository := ci.MapRepository(c)
			assert.Equal(t, debrickedCli, repository)
		})
	}
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
	assert.Equal(t, "", env.Filepath)
	assert.Equal(t, Integration, env.Integration)
	assert.NotEmpty(t, env.Author)
	assert.NotEmpty(t, env.Branch)
	assert.Equal(t, debrickedUrl, env.RepositoryUrl)
	assert.NotEmpty(t, env.Commit)
	assert.Equal(t, debrickedCli, env.Repository)
}
