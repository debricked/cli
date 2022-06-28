package tests

import "testing"

func TestBuildkiteSh(t *testing.T) {
	env := map[string]string{
		"BUILDKITE":        "buildkite",
		"BUILDKITE_COMMIT": validCommit,
		"BUILDKITE_BRANCH": "main",
		"BUILDKITE_REPO":   "https://github.com/debricked/cli.git",
	}
	Test(t, env)
}
