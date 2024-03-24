package file

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
)

const debrickedExclusionEnvVar = "DEBRICKED_EXCLUSIONS"

type DefaultExclusionList struct {
	Directories []string
}

var defaultExclusions = DefaultExclusionList{
	Directories: []string{
		"node_modules",
		"vendor",
		".git",
		"obj",              // nuget
		"bower_components", // bower
	},
}

func DefaultExclusions() []string {
	var exclusions []string
	for _, excluded_dir := range defaultExclusions.Directories {
		exclusions = append(exclusions, filepath.Join("**", excluded_dir, "**"))
	}
	return exclusions
}

func Exclusions() []string {
	values := DefaultExclusions()

	envValue := os.Getenv(debrickedExclusionEnvVar)
	if envValue != "" {
		values = strings.Split(envValue, ",")
	}

	return values
}

func Excluded(exclusions []string, inclusions []string, path string) bool {
	for _, inclusion := range inclusions {
		ex := filepath.Clean(inclusion)
		matched, _ := doublestar.PathMatch(ex, path)
		if matched {
			return false
		}
	}
	for _, exclusion := range exclusions {
		ex := filepath.Clean(exclusion)
		matched, _ := doublestar.PathMatch(ex, path)
		if matched {
			return true
		}
	}

	return false
}
