package resolve

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/debricked/cli/internal/cmd/resolve"
	"github.com/debricked/cli/internal/wire"
	"github.com/stretchr/testify/assert"
)

func removeLines(input, prefix string) string {
	lines := strings.Split(input, "\n")
	var result []string

	for _, line := range lines {
		if !strings.HasPrefix(line, prefix) {
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func TestResolves(t *testing.T) {
	cases := []struct {
		name           string
		manifestFile   string
		lockFileName   string
		expectedFile   string
		packageManager string
	}{
		{
			name:           "basic package.json",
			manifestFile:   "testdata/npm/package.json",
			lockFileName:   "yarn.lock",
			expectedFile:   "testdata/npm/yarn-expected.lock",
			packageManager: "npm",
		},
		{
			name:           "basic requirements.txt",
			manifestFile:   "testdata/pip/requirements.txt",
			lockFileName:   "requirements.txt.pip.debricked.lock",
			expectedFile:   "testdata/pip/expected.lock",
			packageManager: "pip",
		},
		{
			name:           "basic .csproj",
			manifestFile:   "testdata/nuget/csproj/basic.csproj",
			lockFileName:   "packages.lock.json",
			expectedFile:   "testdata/nuget/csproj/packages-expected.lock.json",
			packageManager: "nuget",
		},
		{
			name:           "basic packages.config",
			manifestFile:   "testdata/nuget/packagesconfig/packages.config",
			lockFileName:   "packages.config.nuget.debricked.lock",
			expectedFile:   "testdata/nuget/packagesconfig/packages.config.expected.lock",
			packageManager: "nuget",
		},
	}

	for _, cT := range cases {
		c := cT
		t.Run(c.name, func(t *testing.T) {
			resolveCmd := resolve.NewResolveCmd(wire.GetCliContainer().Resolver())
			lockFileDir := filepath.Dir(c.manifestFile)
			lockFile := filepath.Join(lockFileDir, c.lockFileName)
			// Remove the lock file if it exists
			os.Remove(lockFile)

			err := resolveCmd.RunE(resolveCmd, []string{c.manifestFile})
			assert.NoError(t, err)

			lockFileContents, fileErr := os.ReadFile(lockFile)
			assert.NoError(t, fileErr)

			expectedFileContents, fileErr := os.ReadFile(c.expectedFile)
			assert.NoError(t, fileErr)
			expectedString := string(expectedFileContents)
			actualString := string(lockFileContents)

			if c.packageManager == "pip" {
				// Remove locations as it is dependent on the machine
				expectedString = removeLines(expectedString, "Location: ")
				actualString = removeLines(actualString, "Location: ")

			}
			assert.Equal(t, expectedString, actualString)
		})
	}
}
