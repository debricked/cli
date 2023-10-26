package resolve

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/debricked/cli/internal/cmd/resolve"
	"github.com/debricked/cli/internal/wire"
	"github.com/stretchr/testify/assert"
)

func TestResolves(t *testing.T) {
	cases := []struct {
		name         string
		manifestFile string
		lockFileName string
		expectedFile string
	}{
		{
			name:         "basic package.json",
			manifestFile: "testdata/npm/package.json",
			lockFileName: "yarn.lock",
			expectedFile: "testdata/npm/yarn-expected.lock",
		},
		{
			name:         "basic requirements.txt",
			manifestFile: "testdata/pip/requirements.txt",
			lockFileName: "requirements.txt.debricked.lock",
			expectedFile: "testdata/pip/expected.lock",
		},
		{
			name:         "basic .csproj",
			manifestFile: "testdata/nuget/basic.csproj",
			lockFileName: "packages.lock.json",
			expectedFile: "testdata/nuget/packages-expected.lock.json",
		},
	}

	for _, c := range cases {
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

			assert.Equal(t, string(expectedFileContents), string(lockFileContents))
		})
	}
}
