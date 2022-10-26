package buildkite

import (
	"github.com/debricked/cli/pkg/ci/testdata"
	"os"
	"testing"
)

var buildkiteEnv = map[string]string{
	"BUILDKITE":        "buildkite",
	"BUILDKITE_COMMIT": "commit",
	"BUILDKITE_BRANCH": "main",
	"BUILDKITE_REPO":   "https://github.com/debricked/cli.git",
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
	err := testdata.SetUpCiEnv(buildkiteEnv)
	defer testdata.ResetEnv(buildkiteEnv)
	if err != nil {
		t.Error(err)
	}

	err = testdata.SetUpGitRepository()
	defer testdata.TearDownGitRepository()
	if err != nil {
		t.Error("failed to initialize repository", err)
	}

	ci := Ci{}
	env, err := ci.Map()
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}
	if env.Filepath != "." {
		t.Error("failed to assert that env contained correct filepath")
	}
	if env.Integration != integration {
		t.Error("failed to assert that env contained correct integration")
	}
	if len(env.Author) == 0 {
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

func TestMapRepositoryHttp(t *testing.T) {
	ci := Ci{}
	buildkiteRepo := "https://github.com/debricked/cli.git"
	repository := ci.MapRepository(buildkiteRepo)
	if repository != "debricked/cli" {
		t.Error("failed to assert that repository was set correctly")
	}
	buildkiteRepo = "http://gitlab.com/debricked/cli.git"
	repository = ci.MapRepository(buildkiteRepo)
	if repository != "debricked/cli" {
		t.Error("failed to assert that repository was set correctly")
	}

	buildkiteRepo = "http://scm.com/debricked/cli.git"
	repository = ci.MapRepository(buildkiteRepo)
	if repository != "debricked/cli" {
		t.Error("failed to assert that repository was set correctly")
	}
}

func TestMapRepositoryMisc(t *testing.T) {
	ci := Ci{}
	buildkiteRepo := "git@github.com:debricked/cli.git"
	repository := ci.MapRepository(buildkiteRepo)
	if repository != "debricked/cli" {
		t.Error("failed to assert that repository was set correctly")
	}
	buildkiteRepo = "git@gitlab.com:debricked/cli.git"
	repository = ci.MapRepository(buildkiteRepo)
	if repository != "debricked/cli" {
		t.Error("failed to assert that repository was set correctly")
	}

	buildkiteRepo = "tcp@scm.com:debricked/cli.git"
	repository = ci.MapRepository(buildkiteRepo)
	if repository != "debricked/cli" {
		t.Error("failed to assert that repository was set correctly")
	}
}

func TestMapRepositoryUrlHttp(t *testing.T) {
	ci := Ci{}
	buildkiteRepo := "https://github.com/debricked/cli.git"
	repository := ci.MapRepositoryUrl(buildkiteRepo)
	if repository != "https://github.com/debricked/cli" {
		t.Error("failed to assert that repository was set correctly")
	}
	buildkiteRepo = "http://gitlab.com/debricked/cli.git"
	repository = ci.MapRepositoryUrl(buildkiteRepo)
	if repository != "http://gitlab.com/debricked/cli" {
		t.Error("failed to assert that repository was set correctly")
	}

	buildkiteRepo = "http://gitlab.com/debricked/sub/cli.git"
	repository = ci.MapRepositoryUrl(buildkiteRepo)
	if repository != "http://gitlab.com/debricked/sub/cli" {
		t.Error("failed to assert that repository was set correctly")
	}
}

func TestMapRepositoryUrlGit(t *testing.T) {
	ci := Ci{}
	buildkiteRepo := "git@github.com:debricked/cli.git"
	repository := ci.MapRepositoryUrl(buildkiteRepo)
	if repository != "https://github.com/debricked/cli" {
		t.Error("failed to assert that repository was set correctly")
	}

	buildkiteRepo = "git@gitlab.com:debricked/cli.git"
	repository = ci.MapRepositoryUrl(buildkiteRepo)
	if repository != "https://gitlab.com/debricked/cli" {
		t.Error("failed to assert that repository was set correctly")
	}

	buildkiteRepo = "tcp@scm.com:debricked/sub/cli.git"
	repository = ci.MapRepositoryUrl(buildkiteRepo)
	if repository != "tcp@scm.com:debricked/sub/cli.git" {
		t.Error("failed to assert that repository was set correctly")
	}
}