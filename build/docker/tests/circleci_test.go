package tests

import (
	"testing"
)

func TestCircleciSh(t *testing.T) {
	env := map[string]string{
		"CIRCLECI":                "circleci",
		"CIRCLE_PROJECT_USERNAME": "debricked",
		"CIRCLE_PROJECT_REPONAME": "cli",
		"CIRCLE_SHA1":             validCommit,
		"CIRCLE_BRANCH":           "main",
		"CIRCLE_REPOSITORY_URL":   "https://github.com/debricked/cli.git",
	}
	Test(t, env)

	env["CIRCLE_REPOSITORY_URL"] = "git@github.com:debricked/cli.git"
	Test(t, env)
}
