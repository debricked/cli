package resolve

import (
	"github.com/debricked/cli/pkg/cmd/files/resolve"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestResolvePip(t *testing.T) {
	cases := []struct {
		name             string
		requirementsFile string
		expectedFile     string
	}{
		{
			name:             "basic requirements.txt",
			requirementsFile: "testdata/pip/requirements.txt",
			expectedFile:     "testdata/pip/expected.lock",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resolveCmd := resolve.NewResolveCmd()
			err := resolveCmd.RunE(resolveCmd, []string{c.requirementsFile})
			assert.NoError(t, err)

			lockFileDir := filepath.Dir(c.requirementsFile)
			lockFile := filepath.Join(lockFileDir, ".requirements.txt.debricked.lock")
			lockFileContents, fileErr := os.ReadFile(lockFile)
			assert.NoError(t, fileErr)

			expectedFileContents, fileErr := os.ReadFile(c.expectedFile)
			assert.NoError(t, fileErr)

			assert.Equal(t, string(expectedFileContents), string(lockFileContents))
		})
	}
}
