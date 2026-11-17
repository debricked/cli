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
        "php": "*",
        "ext-json": "*",
        "guzzlehttp/guzzle": "*",
        "acme/auth-client": "*"
    },
    "require-dev": {
        "phpunit/phpunit": "*"
    }
}
`

func TestSanitizeComposerJSON(t *testing.T) {
	sanitized, err := sanitizeComposerJSON([]byte(composerJSONBefore), newTestContext("composer.json"))

	assert.NoError(t, err)
	assert.Equal(t, composerJSONAfter, string(sanitized))
	assert.NotContains(t, string(sanitized), "PROJ-123")
}

func TestSanitizeComposerJSONMissingLists(t *testing.T) {
	sanitized, err := sanitizeComposerJSON([]byte(`{"name":"portal"}`), newTestContext("composer.json"))

	assert.NoError(t, err)
	assert.Equal(t, "{}\n", string(sanitized))
}

func TestSanitizeComposerJSONOmitsEmptySections(t *testing.T) {
	sanitized, err := sanitizeComposerJSON([]byte(`{"require":{"php":">=8.2"},"require-dev":{}}`), newTestContext("composer.json"))

	assert.NoError(t, err)
	assert.Equal(t, "{\n    \"require\": {\n        \"php\": \"*\"\n    }\n}\n", string(sanitized))
}

func TestSanitizeComposerJSONInvalid(t *testing.T) {
	for _, content := range []string{`{not-json`, `{"require": []}`, `{"require-dev": "x"}`} {
		_, err := sanitizeComposerJSON([]byte(content), newTestContext("composer.json"))
		assert.ErrorIs(t, err, util.ErrInvalidFormat, content)
	}
}

func newTestContext(path string) *util.FileContext {
	context, _ := testdata.NewContext(path)

	return context
}
