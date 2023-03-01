package buildkite

import (
	"fmt"
	"os"
	"regexp"

	"github.com/debricked/cli/pkg/ci/env"
	"github.com/debricked/cli/pkg/ci/util"
	"github.com/debricked/cli/pkg/git"
)

const (
	EnvKey      = "BUILDKITE"
	Integration = "buildkite"
)

type Ci struct{}

func (_ Ci) Identify() bool {
	return util.EnvKeyIsSet(EnvKey)
}

func (ci Ci) Map() (env.Env, error) {
	e := env.Env{}
	e.Repository = ci.MapRepository(os.Getenv("BUILDKITE_REPO"))
	e.Commit = os.Getenv("BUILDKITE_COMMIT")
	e.Branch = os.Getenv("BUILDKITE_BRANCH")
	e.RepositoryUrl = ci.MapRepositoryUrl(os.Getenv("BUILDKITE_REPO"))
	e.Integration = Integration
	repo, err := git.FindRepository(e.Filepath)
	if err != nil {

		return e, err
	}
	author, err := git.FindCommitAuthor(repo)
	e.Author = author

	return e, err
}

// MapRepository returns the repository according to the following rules:
//  1. If BUILDKITE_REPO starts with "http(s)://" and ends with ".git", use capture group to find repository.
//  2. If BUILDKITE_REPO starts with "git@" and ends with ".git", use capture group to find repository.
//  3. return BUILDKITE_REPO.
func (_ Ci) MapRepository(buildkiteRepo string) string {
	httpRegex, _ := regexp.Compile(`^https?://.+\.[a-z0-9]+/(.+)\.git$`)
	matches := httpRegex.FindStringSubmatch(buildkiteRepo)
	if len(matches) == 2 {
		return matches[1]
	}

	sshRegex, _ := regexp.Compile(`^.*:[0-9]*/*(.+)\.git$`)
	matches = sshRegex.FindStringSubmatch(buildkiteRepo)
	if len(matches) == 2 {
		return matches[1]
	}

	return buildkiteRepo
}

// MapRepositoryUrl returns the repository url according to the following rules:
//  1. If buildkiteRepo starts with "http(s)://" and ends with ".git", use capture group to find repository.
//  2. If buildkiteRepo is of the form "git@github.com:organisation/reponame.git",
//     rewrite and use "https://github.com/organisation/reponame".
//  3. return buildkiteRepo.
func (_ Ci) MapRepositoryUrl(buildkiteRepo string) string {
	httpRegex, _ := regexp.Compile(`^(https?://.+)\.git$`)
	matches := httpRegex.FindStringSubmatch(buildkiteRepo)
	if len(matches) == 2 {
		return matches[1]
	}

	sshRegex, _ := regexp.Compile(`git@(.+):[0-9]*/?(.+)\.git$`)
	matches = sshRegex.FindStringSubmatch(buildkiteRepo)
	if len(matches) == 3 {
		return fmt.Sprintf("https://%s/%s", matches[1], matches[2])
	}

	return buildkiteRepo
}
