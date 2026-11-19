package composer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/debricked/cli/internal/sanitize/pm/util"
)

const composerLockBefore = `{
  "_readme": ["readme"],
  "content-hash": "aaaa",
  "packages": [
    {
      "name": "acme/auth-client",
      "version": "dev-feature/PROJ-123-new-login",
      "source": {"url": "git@git.acme.internal:platform/auth-client.git"},
      "require": {
        "guzzlehttp/guzzle": "^7.0",
        "php": ">=8.1"
      },
      "authors": [{"name": "Jane Doe"}]
    },
    {
      "name": "acme/shared-utils",
      "version": "1.0.0",
      "dist": {"url": "../shared-utils"}
    }
  ],
  "packages-dev": [
    {
      "name": "phpunit/phpunit",
      "version": "10.5.10",
      "source": {"url": "https://github.com/sebastianbergmann/phpunit.git"}
    }
  ],
  "minimum-stability": "dev"
}
`

const composerLockAfter = `{
    "packages": [
        {
            "name": "acme/auth-client",
            "version": "dev-redacted",
            "require": {
                "guzzlehttp/guzzle": "*",
                "php": "*"
            }
        },
        {
            "name": "acme/shared-utils",
            "version": "1.0.0"
        }
    ],
    "packages-dev": [
        {
            "name": "phpunit/phpunit",
            "version": "10.5.10"
        }
    ]
}
`

func TestSanitizeComposerLock(t *testing.T) {
	sanitized, err := sanitizeComposerLock([]byte(composerLockBefore), newTestContext("composer.lock"))

	assert.NoError(t, err)
	assert.Equal(t, composerLockAfter, string(sanitized))
	assert.NotContains(t, string(sanitized), "PROJ-123")
	assert.NotContains(t, string(sanitized), "acme.internal")
}

func TestSanitizeComposerLockInvalid(t *testing.T) {
	invalid := []string{
		`{not-json`,
		`[]`,
		`{"packages": {}}`,
		`{"packages": ["bad-entry"]}`,
		`{"name":"not-lock"}`,
		`{"packages": [{"version": "1.0.0"}]}`,
		`{"packages": [{"name": "a/b"}]}`,
		`{"packages": [{"name": "a/b", "version": 1}]}`,
		`{"packages": [{"name": "a/b", "version": "1.0.0", "require": ["c/d"]}]}`,
	}

	for _, content := range invalid {
		_, err := sanitizeComposerLock([]byte(content), newTestContext("composer.lock"))
		assert.ErrorIs(t, err, util.ErrInvalidFormat, content)
	}
}

func TestSanitizeComposerLockOmitsEmptyRequire(t *testing.T) {
	before := `{
  "packages": [
    {
      "name": "acme/shared-utils",
      "version": "1.0.0",
      "require": {}
    }
  ]
}
`
	after := `{
    "packages": [
        {
            "name": "acme/shared-utils",
            "version": "1.0.0"
        }
    ]
}
`

	sanitized, err := sanitizeComposerLock([]byte(before), newTestContext("composer.lock"))

	assert.NoError(t, err)
	assert.Equal(t, after, string(sanitized))
}
