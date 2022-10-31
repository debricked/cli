package bitbucket

import (
	"fmt"
	"github.com/debricked/cli/pkg/ci/env"
	"github.com/debricked/cli/pkg/ci/util"
	"github.com/debricked/cli/pkg/git"
	"os"
)

const (
	EnvKey      = "BITBUCKET_BUILD_NUMBER"
	Integration = "bitbucket"
)

type Ci struct{}

func (_ Ci) Identify() bool {
	return util.EnvKeyIsSet(EnvKey)
}

func (_ Ci) Map() (env.Env, error) {
	e := env.Env{}
	e.Repository = fmt.Sprintf("%s/%s", os.Getenv("BITBUCKET_REPO_OWNER"), os.Getenv("BITBUCKET_REPO_SLUG"))
	e.Commit = os.Getenv("BITBUCKET_COMMIT")
	e.Branch = os.Getenv("BITBUCKET_BRANCH")
	e.RepositoryUrl = os.Getenv("BITBUCKET_GIT_HTTP_ORIGIN")
	e.Integration = Integration
	e.Filepath = "."
	repo, err := git.FindRepository(e.Filepath)
	if err != nil {
		return e, nil
	}
	author, err := git.FindCommitAuthor(repo)
	e.Author = author
	return e, err
}
