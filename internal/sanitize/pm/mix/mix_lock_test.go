package mix

import (
	"testing"

	"github.com/debricked/cli/internal/sanitize/pm/util"
	"github.com/stretchr/testify/assert"
)

const mixLockBefore = `%{
  "acme_auth": {:git, "git@git.acme.internal:platform/acme_auth.git", "9f3c1ab", [branch: "feature/PROJ-123"]},
  "acme_logger": {:hex, :acme_logger, "1.4.2", "aaaa", [:mix], [], "hexpm:acme_org", "bbbb"},
  "plug": {:hex, :plug, "1.15.3", "cccc", [:mix], [{:plug_crypto, "~> 2.0", [hex: :plug_crypto, repo: "hexpm", optional: false]}], "hexpm", "dddd"},
  "plug_crypto": {:hex, :plug_crypto, "2.0.0", "eeee", [:mix], [], "hexpm", "ffff"},
}
`

const mixLockAfter = `%{
  "acme_auth": {:git},
  "acme_logger": {:hex, :acme_logger, "1.4.2"},
  "plug": {:hex, :plug, "1.15.3"},
  "plug_crypto": {:hex, :plug_crypto, "2.0.0"},
}
`

func TestSanitizeMixLock(t *testing.T) {
	sanitized, err := sanitizeMixLock([]byte(mixLockBefore), nil)

	assert.NoError(t, err)
	assert.Equal(t, mixLockAfter, string(sanitized))
	assert.NotContains(t, string(sanitized), "acme_org")
}

func TestSanitizeMixLockDropsEntriesOfOtherKinds(t *testing.T) {
	before := "%{\n  \"local\": {:path, \"../local\"},\n  \"plug\": {:hex, :plug, \"1.15.3\"},\n}\n"

	sanitized, err := sanitizeMixLock([]byte(before), nil)

	assert.NoError(t, err)
	assert.Equal(t, "%{\n  \"plug\": {:hex, :plug, \"1.15.3\"},\n}\n", string(sanitized))
}

func TestSanitizeMixLockInvalid(t *testing.T) {
	for _, content := range []string{
		"",
		"%{\n",
		"{\n}\n",
		"%{\n\"plug\": {:hex, :plug, \"1.15.3\"},\n}\n",
		"%{\n  \"plug\": {:hex, \"1.15.3\"},\n}\n",
		"%{\n  plug: 1.15.3\n}\n",
	} {
		_, err := sanitizeMixLock([]byte(content), nil)
		assert.ErrorIs(t, err, util.ErrInvalidFormat, content)
	}
}
