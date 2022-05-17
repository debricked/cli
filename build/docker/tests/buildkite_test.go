package tests

import "testing"

func TestBuildkiteSh(t *testing.T) {
	env := map[string]string{
		"BUILDKITE":        "buildkite",
		"BUILDKITE_COMMIT": "84cac1be9931f8bcc8ef59c5544aaac8c5c97c8b",
		"BUILDKITE_BRANCH": "main",
		"BUILDKITE_REPO":   "https://github.com/debricked/cli.git",
	}
	Test(t, env)
}
