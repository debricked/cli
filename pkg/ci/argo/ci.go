package argo

import (
	"fmt"
	"github.com/debricked/cli/pkg/ci/env"
	"github.com/debricked/cli/pkg/ci/util"
	"github.com/debricked/cli/pkg/git"
	"os"
	"regexp"
)

const (
	EnvKey      = "ARGO_AGENT_TASK_WORKERS"
	Integration = "argoWorkflows"
)

type Ci struct{}

func (_ Ci) Identify() bool {
	return util.EnvKeyIsSet(EnvKey)
}

func (ci Ci) Map() (env.Env, error) {
	e := env.Env{}
	e.Filepath = "."
	e.Repository = ci.MapRepository(os.Getenv("DEBRICKED_GIT_URL"))
	e.RepositoryUrl = ci.MapRepositoryUrl(os.Getenv("DEBRICKED_GIT_URL"))
	e.Integration = Integration
	repo, err := git.FindRepository(e.Filepath)
	if err != nil {
		return e, err
	}

	var gitErrs []error

	commit, commitErr := git.FindCommitHash(repo)
	if commitErr != nil {
		gitErrs = append(gitErrs, commitErr)
	}
	e.Commit = commit

	branch, branchErr := git.FindBranch(repo)
	if branchErr != nil {
		gitErrs = append(gitErrs, branchErr)
	}
	e.Branch = branch

	author, authorErr := git.FindCommitAuthor(repo)
	if authorErr != nil {
		gitErrs = append(gitErrs, authorErr)
	}
	e.Author = author

	if len(gitErrs) > 0 {
		err = gitErrs[0]
	}

	return e, err
}

// MapRepository returns repository according to the following rules:
//  1. If gitUrl starts with "http(s)://" and ends with ".git", use capture group to set repository.
//  2. If gitUrl starts with "git@" and ends with ".git", use capture group to set repository.
//  3. Return gitUrl.
func (_ Ci) MapRepository(gitUrl string) string {
	httpRegex, _ := regexp.Compile(`^https?://.+\.[a-z0-9]+/(.+)\.git$`)
	matches := httpRegex.FindStringSubmatch(gitUrl)
	if len(matches) == 2 {
		return matches[1]
	}

	sshRegex, _ := regexp.Compile(`^.*:[0-9]*/*(.+)\.git$`)
	matches = sshRegex.FindStringSubmatch(gitUrl)
	if len(matches) == 2 {
		return matches[1]
	}

	return gitUrl
}

// MapRepositoryUrl returns repository URL according to the following rules:
//  1. If gitUrl starts with "http(s)://" and ends with ".git", use capture group to set repository URL.
//  2. If gitUrl is of the form "git@github.com:organisation/reponame.git",
//     rewrite and use "https://github.com/organisation/reponame" as repository URL.
//  3. Otherwise, return gitUrl
func (_ Ci) MapRepositoryUrl(gitUrl string) string {
	httpRegex, _ := regexp.Compile(`^(https?://.+)\.git$`)
	matches := httpRegex.FindStringSubmatch(gitUrl)
	if len(matches) == 2 {
		return matches[1]
	}

	sshRegex, _ := regexp.Compile(`git@(.+):[0-9]*/?(.+)\.git$`)
	matches = sshRegex.FindStringSubmatch(gitUrl)
	if len(matches) == 3 {
		return fmt.Sprintf("https://%s/%s", matches[1], matches[2])
	}

	return gitUrl
}
