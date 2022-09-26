package travis

import (
	"debricked/pkg/ci/env"
	"debricked/pkg/ci/util"
	"debricked/pkg/git"
	"fmt"
	"os"
)

const (
	EnvKey      = "TRAVIS"
	integration = "travis"
)

type Ci struct{}

func (_ Ci) Identify() bool {
	return util.EnvKeyIsSet(EnvKey)
}

func (_ Ci) Map() (env.Env, error) {
	e := env.Env{}
	e.Repository = os.Getenv("TRAVIS_REPO_SLUG")
	e.Commit = os.Getenv("TRAVIS_COMMIT")
	e.Branch = os.Getenv("TRAVIS_BRANCH")
	e.RepositoryUrl = fmt.Sprintf("https://github.com/%s", e.Repository)
	e.Integration = integration
	//# The absolute path to the directory where the repository being built has been copied on the worker.
	//# HOME is set to /home/travis on Linux, /Users/travis on MacOS, and /c/Users/travis on Windows.
	e.Filepath = os.Getenv("TRAVIS_BUILD_DIR")
	repo, err := git.FindRepository(e.Filepath)
	if err != nil {
		return e, nil
	}
	author, err := git.FindCommitAuthor(repo)
	e.Author = author
	return e, err
}
