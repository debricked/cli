package tests

import (
	"testing"
)

func TestGithubSh(t *testing.T) {
	env := map[string]string{
		"GITHUB_ACTION":     "githubActions",
		"GITHUB_REPOSITORY": "debricked/cli",
		"GITHUB_SHA":        "84cac1be9931f8bcc8ef59c5544aaac8c5c97c8b",
		"GITHUB_REF":        "refs/heads/main",
	}
	Test(t, env)
}
