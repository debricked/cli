package gitlab

import (
	"testing"

	"github.com/debricked/cli/internal/ci/testdata"
	"github.com/stretchr/testify/assert"
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
	testdata.AssertIdentify(t, Ci{}.Identify, EnvKey)
}

func TestParse(t *testing.T) {
	testdata.SetUpCiEnv(t, gitLabEnv)
	defer testdata.ResetEnv(t, gitLabEnv)
	ci := Ci{}

	env, _ := ci.Map()

	assert.Equal(t, gitLabEnv["CI_PROJECT_DIR"], env.Filepath)
	assert.Equal(t, Integration, env.Integration)
	assert.Equal(t, gitLabEnv["CI_COMMIT_AUTHOR"], env.Author)
	assert.Equal(t, gitLabEnv["CI_COMMIT_REF_NAME"], env.Branch)
	assert.Equal(t, "https://gitlab.com/debricked/cli", env.RepositoryUrl)
	assert.Equal(t, gitLabEnv["CI_COMMIT_SHA"], env.Commit)
	assert.Equal(t, "debricked/cli", env.Repository)
}
