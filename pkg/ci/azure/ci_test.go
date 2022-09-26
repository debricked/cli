package azure

import (
	"debricked/pkg/ci/testdata"
	"os"
	"testing"
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
	err := testdata.SetUpCiEnv(azureEnv)
	if err != nil {
		t.Fatal(err)
	}
	defer testdata.ResetEnv(azureEnv)

	ci := Ci{}
	env, _ := ci.Map()
	if env.Filepath != azureEnv["BUILD_SOURCESDIRECTORY"] {
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
