package circleci

import (
	"fmt"
	"github.com/debricked/cli/pkg/ci/env"
	"github.com/debricked/cli/pkg/ci/util"
	"github.com/debricked/cli/pkg/git"
	"os"
	"regexp"
)

const (
	EnvKey      = "CIRCLECI"
	integration = "circleci"
)

type Ci struct{}

func (_ Ci) Identify() bool {
	return util.EnvKeyIsSet(EnvKey)
}

func (ci Ci) Map() (env.Env, error) {
	e := env.Env{}
	e.Repository = fmt.Sprintf("%s/%s", os.Getenv("CIRCLE_PROJECT_USERNAME"), os.Getenv("CIRCLE_PROJECT_REPONAME"))
	e.Commit = os.Getenv("CIRCLE_SHA1")
	e.Branch = os.Getenv("CIRCLE_BRANCH")
	e.RepositoryUrl = ci.MapRepositoryUrl(os.Getenv("CIRCLE_REPOSITORY_URL"))
	e.Integration = integration
	e.Filepath = "."
	repo, err := git.FindRepository(e.Filepath)
	if err != nil {
		return e, nil
	}
	author, err := git.FindCommitAuthor(repo)
	e.Author = author
	return e, err
}

// MapRepositoryUrl returns the repository url according to the following rules:
//  1. If circleCiRepo starts with "http(s)://", use it as the repo url.
//  2. If circleCiRepo is of the form "git@github.com:organisation/reponame.git",
//     rewrite and use "https://github.com/organisation/reponame" as repo url.
//  3. return circleCiRepo
func (_ Ci) MapRepositoryUrl(circleCiRepo string) string {
	httpRegex, _ := regexp.Compile("^(https?://.+)\\.git$")
	matches := httpRegex.FindStringSubmatch(circleCiRepo)
	if len(matches) == 2 {
		return matches[1]
	}

	sshRegex, _ := regexp.Compile("git@(.+):[0-9]*/?(.+)\\.git$")
	matches = sshRegex.FindStringSubmatch(circleCiRepo)
	if len(matches) == 3 {
		return fmt.Sprintf("https://%s/%s", matches[1], matches[2])
	}

	return circleCiRepo
}
