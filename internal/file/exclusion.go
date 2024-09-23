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
		".vscode-test",     // excluding testing framework
	},
}

func DefaultExclusions() []string {
	var exclusions []string
	for _, excluded_dir := range defaultExclusions.Directories {
		exclusions = append(exclusions, "**/"+excluded_dir+"/**")
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
	path = filepath.ToSlash(path)
	for _, inclusion := range inclusions {
		matched, _ := doublestar.Match(inclusion, path)
		if matched {
			return false
		}
	}
	for _, exclusion := range exclusions {
		matched, _ := doublestar.Match(exclusion, path)
		if matched {
			return true
		}
	}

	return false
}
