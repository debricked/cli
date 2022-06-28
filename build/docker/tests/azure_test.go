package tests

import (
	"testing"
)

func TestAzureSh(t *testing.T) {
	env := map[string]string{
		"TF_BUILD":               "azure",
		"SYSTEM_COLLECTIONURI":   "debricked",
		"BUILD_REPOSITORY_NAME":  "cli",
		"BUILD_SOURCEVERSION":    validCommit,
		"BUILD_SOURCEBRANCHNAME": "main",
		"BUILD_SOURCESDIRECTORY": ".",
		"BUILD_REPOSITORY_URI":   "https://github.com/debricked/cli",
	}
	Test(t, env)
}
