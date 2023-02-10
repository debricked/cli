package argo

import (
	"github.com/debricked/cli/pkg/ci/env"
	"github.com/debricked/cli/pkg/ci/testdata"
	"os"
	"testing"
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
	err := testdata.SetUpCiEnv(argoEnv)
	if err != nil {
		t.Error(err)
	}
	defer testdata.ResetEnv(argoEnv, t)

	cwd, err := testdata.SetUpGitRepository(true)
	if err != nil {
		t.Fatal("failed to initialize repository", err)
	}
	defer testdata.TearDownGitRepository(cwd, t)

	ci := Ci{}
	env, err := ci.Map()
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}
	assertEnv(env, t)
}

func TestMapRepositoryHttp(t *testing.T) {
	ci := Ci{}
	buildkiteRepo := "https://github.com/debricked/cli.git"
	repository := ci.MapRepository(buildkiteRepo)
	if repository != debrickedCli {
		t.Error("failed to assert that repository was set correctly")
	}
	buildkiteRepo = "http://gitlab.com/debricked/cli.git"
	repository = ci.MapRepository(buildkiteRepo)
	if repository != debrickedCli {
		t.Error("failed to assert that repository was set correctly")
	}

	buildkiteRepo = "http://scm.com/debricked/cli.git"
	repository = ci.MapRepository(buildkiteRepo)
	if repository != debrickedCli {
		t.Error("failed to assert that repository was set correctly")
	}
}

func TestMapRepositoryMisc(t *testing.T) {
	ci := Ci{}
	buildkiteRepo := "git@github.com:debricked/cli.git"
	repository := ci.MapRepository(buildkiteRepo)
	if repository != debrickedCli {
		t.Error("failed to assert that repository was set correctly")
	}
	buildkiteRepo = "git@gitlab.com:debricked/cli.git"
	repository = ci.MapRepository(buildkiteRepo)
	if repository != debrickedCli {
		t.Error("failed to assert that repository was set correctly")
	}

	buildkiteRepo = "tcp@scm.com:debricked/cli.git"
	repository = ci.MapRepository(buildkiteRepo)
	if repository != debrickedCli {
		t.Error("failed to assert that repository was set correctly")
	}
}

func TestMapRepositoryUrlHttp(t *testing.T) {
	ci := Ci{}
	buildkiteRepo := "https://github.com/debricked/cli.git"
	repository := ci.MapRepositoryUrl(buildkiteRepo)
	if repository != debrickedUrl {
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
	if repository != debrickedUrl {
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

func assertEnv(env env.Env, t *testing.T) {
	if env.Filepath != "" {
		t.Error("failed to assert that env contained correct filepath")
	}
	if env.Integration != Integration {
		t.Error("failed to assert that env contained correct integration")
	}
	if len(env.Author) == 0 {
		t.Error("failed to assert that env contained correct author")
	}
	if len(env.Branch) == 0 {
		t.Error("failed to assert that env contained correct branch")
	}
	if env.RepositoryUrl != debrickedUrl {
		t.Error("failed to assert that env contained correct repository URL")
	}
	if len(env.Commit) == 0 {
		t.Error("failed to assert that env contained correct commit")
	}
	if env.Repository != debrickedCli {
		t.Error("failed to assert that env contained correct repository")
	}
}
