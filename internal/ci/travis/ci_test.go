package travis

import (
	"testing"

	"github.com/debricked/cli/internal/ci/env"
	"github.com/debricked/cli/internal/ci/testdata"
	"github.com/stretchr/testify/assert"
)

var travisEnv = map[string]string{
	"TRAVIS_REPO_SLUG": "debricked/cli",
	"TRAVIS_BRANCH":    "main",
	"TRAVIS_COMMIT":    "commit",
	"TRAVIS_BUILD_DIR": ".",
}

func TestIdentify(t *testing.T) {
	testdata.AssertIdentify(t, Ci{}.Identify, EnvKey)
}

func TestParse(t *testing.T) {
	testdata.SetUpCiEnv(t, travisEnv)
	defer testdata.ResetEnv(t, travisEnv)

	cwd := testdata.SetUpGitRepository(t, true)
	defer testdata.TearDownGitRepository(cwd, t)

	ci := Ci{}
	e, err := ci.Map()
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}
	assertEnv(t, e)
}

func assertEnv(t *testing.T, env env.Env) {
	assert.Equal(t, travisEnv["TRAVIS_BUILD_DIR"], env.Filepath)
	assert.Equal(t, Integration, env.Integration)
	assert.NotEmpty(t, env.Author)
	assert.Equal(t, travisEnv["TRAVIS_BRANCH"], env.Branch)
	assert.Equal(t, "https://github.com/debricked/cli", env.RepositoryUrl)
	assert.Equal(t, travisEnv["TRAVIS_COMMIT"], env.Commit)
	assert.Equal(t, "debricked/cli", env.Repository)
}
