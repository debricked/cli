package github

import (
	"fmt"
	"github.com/debricked/cli/pkg/ci/env"
	"github.com/debricked/cli/pkg/ci/util"
	"os"
)

const (
	EnvKey      = "GITHUB_ACTION"
	Integration = "githubActions"
)

type Ci struct{}

func (_ Ci) Identify() bool {
	return util.EnvKeyIsSet(EnvKey)
}

func (_ Ci) Map() (env.Env, error) {
	e := env.Env{}
	e.Repository = os.Getenv("GITHUB_REPOSITORY")
	e.Commit = os.Getenv("GITHUB_SHA")
	e.Branch = os.Getenv("GITHUB_REF_NAME")
	e.RepositoryUrl = fmt.Sprintf("https://github.com/%s", os.Getenv("GITHUB_REPOSITORY"))
	e.Integration = Integration
	fmt.Println("Integration:", e.Integration)
	e.Filepath = "."
	e.Author = os.Getenv("GITHUB_ACTOR")
	return e, nil
}
