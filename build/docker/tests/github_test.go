package tests

import (
	"testing"
)

func TestGithubSh(t *testing.T) {
	env := map[string]string{
		"GITHUB_ACTION":     "githubActions",
		"GITHUB_REPOSITORY": "debricked/cli",
		"GITHUB_SHA":        validCommit,
		"GITHUB_REF":        "refs/heads/main",
	}
	Test(t, env)
}
