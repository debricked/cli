package tests

import "testing"

func TestBitbucketSh(t *testing.T) {
	env := map[string]string{
		"BITBUCKET_BUILD_NUMBER":    "2",
		"BITBUCKET_REPO_OWNER":      "debricked",
		"BITBUCKET_REPO_SLUG":       "cli",
		"BITBUCKET_COMMIT":          validCommit,
		"BITBUCKET_BRANCH":          "main",
		"BITBUCKET_GIT_HTTP_ORIGIN": "https://github.com/debricked/cli",
	}
	Test(t, env)
}
