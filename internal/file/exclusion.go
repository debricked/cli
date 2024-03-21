package file

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

const debrickedExclusionEnvVar = "DEBRICKED_EXCLUSIONS"

func DefaultExclusions() []string {
	return []string{
		filepath.Join("**", "node_modules", "**"),
		filepath.Join("**", "vendor", "**"),
		filepath.Join("**", ".git", "**"),
		filepath.Join("**", "obj", "**"),              // nuget
		filepath.Join("**", "bower_components", "**"), // bower
	}
}

func Exclusions() []string {
	values := DefaultExclusions()

	envValue := os.Getenv(debrickedExclusionEnvVar)
	if envValue != "" {
		values = strings.Split(envValue, ",")
	}

	return values
}

var EXCLUDED_DIRS_FINGERPRINT = []string{
	"nbproject", "nbbuild", "nbdist", "node_modules",
	"__pycache__", "_yardoc", "eggs",
	"wheels", "htmlcov", "__pypackages__", ".git"}

var EXCLUDED_DIRS_FINGERPRINT_RAW = []string{"**/*.egg-info/**", "**/*venv/**", "**/*venv3/**"}

func DefaultExclusionsFingerprint() []string {
	output := []string{}

	for _, pattern := range EXCLUDED_DIRS_FINGERPRINT {
		output = append(output, filepath.Join("**", pattern, "**"))
	}

	output = append(output, EXCLUDED_DIRS_FINGERPRINT_RAW...)

	return output
}

func Excluded(exclusions []string, path string) bool {
	for _, exclusion := range exclusions {
		ex := filepath.Clean(exclusion)
		matched, _ := doublestar.PathMatch(ex, path)
		if matched {
			return true
		}
	}

	return false
}
