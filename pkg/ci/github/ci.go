package github

import (
	"fmt"
	"github.com/debricked/cli/pkg/ci/env"
	"github.com/debricked/cli/pkg/ci/util"
	"os"
	"strings"
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

	// Github gives branches as: refs/heads/master, and tags as refs/tags/v1.1.0.
	// Remove prefix refs/{tags,heads} from the name before sending to Debricked.
	gitHubRef := os.Getenv("GITHUB_REF")
	branch := strings.Replace(gitHubRef, "refs/heads/", "", 1)
	branch = strings.Replace(branch, "refs/tags/", "", 1)
	if strings.Contains(branch, "/merge") {
		branch = os.Getenv("GITHUB_HEAD_REF")
	}
	e.Branch = branch

	e.RepositoryUrl = fmt.Sprintf("https://github.com/%s", os.Getenv("GITHUB_REPOSITORY"))
	e.Integration = Integration
	e.Filepath = "."
	e.Author = os.Getenv("GITHUB_ACTOR")
	return e, nil
}
