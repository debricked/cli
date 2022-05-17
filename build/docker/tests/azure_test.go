package tests

import (
	"testing"
)

func TestAzureSh(t *testing.T) {
	env := map[string]string{
		"TF_BUILD":               "azure",
		"SYSTEM_COLLECTIONURI":   "debricked",
		"BUILD_REPOSITORY_NAME":  "cli",
		"BUILD_SOURCEVERSION":    "84cac1be9931f8bcc8ef59c5544aaac8c5c97c8b",
		"BUILD_SOURCEBRANCHNAME": "main",
		"BUILD_SOURCESDIRECTORY": ".",
		"BUILD_REPOSITORY_URI":   "https://github.com/debricked/cli",
	}
	Test(t, env)
}
