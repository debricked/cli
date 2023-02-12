package bitbucket

import (
	"github.com/debricked/cli/pkg/ci/env"
	"github.com/debricked/cli/pkg/ci/testdata"
	"github.com/stretchr/testify/assert"
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
	testdata.AssertIdentify(t, Ci{}.Identify, EnvKey)
}

func TestParse(t *testing.T) {
	testdata.SetUpCiEnv(t, bitbucketEnv)
	defer testdata.ResetEnv(t, bitbucketEnv)

	cwd := testdata.SetUpGitRepository(t, true)
	defer testdata.TearDownGitRepository(cwd, t)

	ci := Ci{}
	e, err := ci.Map()
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}
	assertEnv(e, t)
}

func assertEnv(env env.Env, t *testing.T) {
	assert.Empty(t, env.Filepath)
	assert.Equal(t, Integration, env.Integration)
	assert.NotEmpty(t, env.Author)
	assert.Equal(t, bitbucketEnv["BITBUCKET_BRANCH"], env.Branch)
	assert.Equal(t, bitbucketEnv["BITBUCKET_GIT_HTTP_ORIGIN"], env.RepositoryUrl)
	assert.Equal(t, bitbucketEnv["BITBUCKET_COMMIT"], env.Commit)
	assert.Equal(t, "debricked/cli", env.Repository)
}
