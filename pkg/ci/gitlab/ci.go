package gitlab

import (
	"github.com/debricked/cli/pkg/ci/env"
	"github.com/debricked/cli/pkg/ci/util"
	"os"
)

const (
	EnvKey      = "GITLAB_CI"
	integration = "gitlab"
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
	e.Integration = integration
	e.Filepath = os.Getenv("CI_PROJECT_DIR")
	e.Author = os.Getenv("CI_COMMIT_AUTHOR")
	return e, nil
}
