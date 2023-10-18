package file

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultExclusions(t *testing.T) {
	separator := string(os.PathSeparator)
	for _, ex := range DefaultExclusions() {
		exParts := strings.Split(ex, separator)
		assert.Greaterf(t, len(exParts), 0, "failed to assert that %s used correct separator. Proper separator %s", ex, separator)
	}
}

func TestDefaultExclusionsFingerprint(t *testing.T) {
	expectedExclusions := []string{
		filepath.Join("**", "nbproject", "**"),
		filepath.Join("**", "nbbuild", "**"),
		filepath.Join("**", "nbdist", "**"),
		filepath.Join("**", "node_modules", "**"),
		filepath.Join("**", "__pycache__", "**"),
		filepath.Join("**", "_yardoc", "**"),
		filepath.Join("**", "eggs", "**"),
		filepath.Join("**", "wheels", "**"),
		filepath.Join("**", "htmlcov", "**"),
		filepath.Join("**", "__pypackages__", "**"),
		"**/*.egg-info/**",
		"**/*venv/**",
	}

	exclusions := DefaultExclusionsFingerprint()

	assert.ElementsMatch(t, expectedExclusions, exclusions, "DefaultExclusionsFingerprint did not return the expected exclusions")
}

func TestDefaultExclusionsWithTokenEnvVariable(t *testing.T) {
	oldEnvValue := os.Getenv(debrickedExclusionEnvVar)
	gt := []string{"*/**.lock", "**/node_modules/**", "*\\**.ex"}
	err := os.Setenv(debrickedExclusionEnvVar, "*/**.lock,**/node_modules/**,*\\**.ex")

	if err != nil {
		t.Fatalf("failed to set env var %s", debrickedExclusionEnvVar)
	}

	defer func(key, value string) {
		err := os.Setenv(key, value)
		if err != nil {
			t.Fatalf("failed to reset env var %s", debrickedExclusionEnvVar)
		}
	}(debrickedExclusionEnvVar, oldEnvValue)

	defExclusions := DefaultExclusions()
	assert.Equal(t, gt, defExclusions)

}

func TestDefaultExclusionsWithEmptyTokenEnvVariable(t *testing.T) {
	oldEnvValue := os.Getenv(debrickedExclusionEnvVar)
	err := os.Setenv(debrickedExclusionEnvVar, "")
	gt := []string{
		filepath.Join("**", "node_modules", "**"),
		filepath.Join("**", "vendor", "**"),
		filepath.Join("**", ".git", "**"),
		filepath.Join("**", "obj", "**"),
	}

	if err != nil {
		t.Fatalf("failed to set env var %s", debrickedExclusionEnvVar)
	}

	defer func(key, value string) {
		err := os.Setenv(key, value)
		if err != nil {
			t.Fatalf("failed to reset env var %s", debrickedExclusionEnvVar)
		}
	}(debrickedExclusionEnvVar, oldEnvValue)

	defExclusions := DefaultExclusions()
	assert.Equal(t, gt, defExclusions)
}
