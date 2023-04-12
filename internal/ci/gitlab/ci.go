package gitlab

import (
	"os"

	"github.com/debricked/cli/internal/ci/env"
	"github.com/debricked/cli/internal/ci/util"
)

const (
	EnvKey      = "GITLAB_CI"
	Integration = "gitlab"
)

type Ci struct{}

func (_ Ci) Identify() bool {
	return util.EnvKeyIsSet(EnvKey)
}

func (_ Ci) Map() (env.Env, error) {
	e := env.Env{}
	e.Repository = os.Getenv("CI_PROJECT_PATH")
	e.Commit = os.Getenv("CI_COMMIT_SHA")
	e.Branch = os.Getenv("CI_COMMIT_REF_NAME")
	e.RepositoryUrl = os.Getenv("CI_PROJECT_URL")
	e.Integration = Integration
	e.Filepath = os.Getenv("CI_PROJECT_DIR")
	e.Author = os.Getenv("CI_COMMIT_AUTHOR")

	return e, nil
}
