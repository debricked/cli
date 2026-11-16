package composer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/debricked/cli/internal/sanitize/pm/testdata"
	"github.com/debricked/cli/internal/sanitize/pm/util"
)

const composerJSONBefore = `{
  "name": "acme/payments-portal",
  "description": "Internal card payments portal",
  "type": "project",
  "repositories": [
    {"type": "composer", "url": "https://satis.acme.internal"}
  ],
  "require": {
    "php": ">=8.2",
    "ext-json": "*",
    "guzzlehttp/guzzle": "^7.8",
    "acme/auth-client": "dev-feature/PROJ-123-new-login"
  },
  "require-dev": {
    "phpunit/phpunit": "^10.5"
  },
  "scripts": {
    "deploy": "s3cret"
  }
}
`

const composerJSONAfter = `{
    "require": {
    "php": ">=8.2",
    "ext-json": "*",
    "guzzlehttp/guzzle": "^7.8",
    "acme/auth-client": "dev-feature/PROJ-123-new-login"
    },
    "require-dev": {
    "phpunit/phpunit": "^10.5"
    }
}
`

func TestSanitizeComposerJSON(t *testing.T) {
	sanitized, err := sanitizeComposerJSON([]byte(composerJSONBefore), newTestContext("composer.json"))

	assert.NoError(t, err)
  assert.JSONEq(t, composerJSONAfter, string(sanitized))
}

func TestSanitizeComposerJSONMissingLists(t *testing.T) {
	sanitized, err := sanitizeComposerJSON([]byte(`{"name":"portal"}`), newTestContext("composer.json"))

	assert.NoError(t, err)
	assert.Equal(t, "{}\n", string(sanitized))
}

func TestSanitizeComposerJSONOmitsEmptySections(t *testing.T) {
	sanitized, err := sanitizeComposerJSON([]byte(`{"require":{"php":">=8.2"},"require-dev":{}}`), newTestContext("composer.json"))

	assert.NoError(t, err)
  assert.JSONEq(t, `{
    "require": {
    "php": ">=8.2"
    }
}`, string(sanitized))
}

func TestSanitizeComposerJSONInvalid(t *testing.T) {
	_, err := sanitizeComposerJSON([]byte(`{not-json`), newTestContext("composer.json"))

	assert.ErrorIs(t, err, util.ErrInvalidFormat)
}

func newTestContext(path string) *util.FileContext {
	context, _ := testdata.NewContext(path)

	return context
}
