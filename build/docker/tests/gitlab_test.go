package tests

import (
	"testing"
)

func TestGitlabSh(t *testing.T) {
	env := map[string]string{
		"GITLAB_CI":          "gitlab",
		"CI_PROJECT_PATH":    "debricked/cli",
		"CI_COMMIT_SHA":      validCommit,
		"CI_COMMIT_REF_NAME": "main",
		"CI_PROJECT_DIR":     "/",
		"CI_PROJECT_URL":     "https://gitlab.com/debricked/cli",
	}
	Test(t, env)
}
