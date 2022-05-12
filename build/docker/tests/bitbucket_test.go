package tests

import "testing"

func TestBitbucketSh(t *testing.T) {
	env := map[string]string{
		"BITBUCKET_BUILD_NUMBER":    "2",
		"BITBUCKET_REPO_OWNER":      "debricked",
		"BITBUCKET_REPO_SLUG":       "cli",
		"BITBUCKET_COMMIT":          "84cac1be9931f8bcc8ef59c5544aaac8c5c97c8b",
		"BITBUCKET_BRANCH":          "main",
		"BITBUCKET_GIT_HTTP_ORIGIN": "https://github.com/debricked/cli",
	}
	Test(t, env)
}
