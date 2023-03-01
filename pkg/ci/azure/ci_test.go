package azure

import (
	"testing"

	"github.com/debricked/cli/pkg/ci/testdata"
	"github.com/stretchr/testify/assert"
)

var azureEnv = map[string]string{
	"TF_BUILD":                "azure",
	"SYSTEM_COLLECTIONURI":    "dir/debricked/",
	"BUILD_REPOSITORY_NAME":   "cli",
	"BUILD_SOURCEVERSION":     "commit",
	"BUILD_SOURCEBRANCHNAME":  "main",
	"BUILD_SOURCESDIRECTORY":  ".",
	"BUILD_REPOSITORY_URI":    "https://github.com/debricked/cli",
	"BUILD_REQUESTEDFOREMAIL": "viktigpetterr <test@test.com>",
}

func TestIdentify(t *testing.T) {
	testdata.AssertIdentify(t, Ci{}.Identify, EnvKey)
}

func TestParse(t *testing.T) {
	testdata.SetUpCiEnv(t, azureEnv)
	defer testdata.ResetEnv(t, azureEnv)
	ci := Ci{}

	env, _ := ci.Map()

	assert.Equal(t, azureEnv["BUILD_SOURCESDIRECTORY"], env.Filepath)
	assert.Equal(t, Integration, env.Integration)
	assert.Equal(t, azureEnv["BUILD_REQUESTEDFOREMAIL"], env.Author)
	assert.Equal(t, azureEnv["BUILD_SOURCEBRANCHNAME"], env.Branch)
	assert.Equal(t, azureEnv["BUILD_REPOSITORY_URI"], env.RepositoryUrl)
	assert.Equal(t, azureEnv["BUILD_SOURCEVERSION"], env.Commit)
	assert.Equal(t, "debricked/cli", env.Repository)

}
