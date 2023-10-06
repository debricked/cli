package file

import "path/filepath"

func DefaultExclusions() []string {
	return []string{
		filepath.Join("**", "node_modules", "**"),
		filepath.Join("**", "vendor", "**"),
		filepath.Join("**", ".git", "**"),
		filepath.Join("**", "obj", "**"), // nuget
	}
}

var EXCLUDED_DIRS_FINGERPRINT = []string{"nbproject", "nbbuild", "nbdist",
	"__pycache__", "venv", "_yardoc", "eggs",
	"wheels", "htmlcov", "__pypackages__", ".egg-info"}

func DefaultExclusionsFingerprint() []string {
	output := []string{}

	for _, pattern := range EXCLUDED_DIRS_FINGERPRINT {
		output = append(output, filepath.Join("**", pattern, "**"))
	}

	return output
}
