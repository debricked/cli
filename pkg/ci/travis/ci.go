package travis

import (
	"fmt"
	"os"

	"github.com/debricked/cli/pkg/ci/env"
	"github.com/debricked/cli/pkg/ci/util"
	"github.com/debricked/cli/pkg/git"
)

const (
	EnvKey      = "TRAVIS"
	Integration = "travis"
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
	e.Integration = Integration
	//# The absolute path to the directory where the repository being built has been copied on the worker.
	//# HOME is set to /home/travis on Linux, /Users/travis on MacOS, and /c/Users/travis on Windows.
	e.Filepath = os.Getenv("TRAVIS_BUILD_DIR")
	repo, err := git.FindRepository(e.Filepath)
	if err != nil {
		return e, err
	}
	author, err := git.FindCommitAuthor(repo)
	e.Author = author

	return e, err
}
