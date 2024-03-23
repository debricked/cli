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

func TestExclusionsWithTokenEnvVariable(t *testing.T) {
	oldEnvValue := os.Getenv(debrickedExclusionEnvVar)
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

	gt := []string{"*/**.lock", "**/node_modules/**", "*\\**.ex"}
	exclusions := Exclusions()
	assert.Equal(t, gt, exclusions)

}

func TestExclusionsWithEmptyTokenEnvVariable(t *testing.T) {
	oldEnvValue := os.Getenv(debrickedExclusionEnvVar)
	err := os.Setenv(debrickedExclusionEnvVar, "")

	if err != nil {
		t.Fatalf("failed to set env var %s", debrickedExclusionEnvVar)
	}

	defer func(key, value string) {
		err := os.Setenv(key, value)
		if err != nil {
			t.Fatalf("failed to reset env var %s", debrickedExclusionEnvVar)
		}
	}(debrickedExclusionEnvVar, oldEnvValue)

	gt := []string{
		filepath.Join("**", "node_modules", "**"),
		filepath.Join("**", "vendor", "**"),
		filepath.Join("**", ".git", "**"),
		filepath.Join("**", "obj", "**"),
		filepath.Join("**", "bower_components", "**"),
	}
	defaultExclusions := Exclusions()
	assert.Equal(t, gt, defaultExclusions)
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
		filepath.Join("**", ".git", "**"),
		"**/*.egg-info/**",
		"**/*venv/**",
		"**/*venv3/**",
	}

	exclusions := DefaultExclusionsFingerprint()

	assert.ElementsMatch(t, expectedExclusions, exclusions, "DefaultExclusionsFingerprint did not return the expected exclusions")
}

func TestExclude(t *testing.T) {
	var files []string
	_ = filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			files = append(files, path)

			return nil
		})

	cases := []struct {
		name               string
		exclusions         []string
		expectedExclusions []string
	}{
		{
			name:               "NoExclusions",
			exclusions:         []string{},
			expectedExclusions: []string{},
		},
		{
			name:               "InvalidFileExclusion",
			exclusions:         []string{"composer.json"},
			expectedExclusions: []string{},
		},
		{
			name:               "FileExclusionWithDoublestar",
			exclusions:         []string{"**/composer.json"},
			expectedExclusions: []string{"composer.json", "composer.json"}, // Two composer.json files in testdata folder
		},
		{
			name:               "DirectoryExclusion",
			exclusions:         []string{"*/composer/*"},
			expectedExclusions: []string{"composer.json", "composer.lock"},
		},
		{
			name:               "DirectoryExclusionWithRelPath",
			exclusions:         []string{"testdata/go/*"},
			expectedExclusions: []string{"go.mod"},
		},
		{
			name:               "ExtensionExclusionWithWildcardAndDoublestar",
			exclusions:         []string{"**/*.mod"},
			expectedExclusions: []string{"go.mod", "go.mod"}, // Two go.mod files in testdata folder
		},
		{
			name:               "DirectoryExclusionWithDoublestar",
			exclusions:         []string{"**/yarn/**"},
			expectedExclusions: []string{"yarn", "yarn.lock"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var excludedFiles []string
			for _, file := range files {
				if Excluded(c.exclusions, []string{}, file) {
					excludedFiles = append(excludedFiles, file)
				}
			}

			assert.Equal(t, len(c.expectedExclusions), len(excludedFiles), "failed to assert that the same number of files were ignored")

			for _, file := range excludedFiles {
				baseName := filepath.Base(file)
				asserted := false
				for _, expectedExcludedFile := range c.expectedExclusions {
					if baseName == expectedExcludedFile {
						asserted = true

						break
					}
				}

				assert.Truef(t, asserted, "%s ignored when it should pass", file)
			}
		})
	}
}
