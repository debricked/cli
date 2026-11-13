package bundler

import (
	"strings"
	"testing"

	"github.com/debricked/cli/internal/sanitize/pm/util"
	"github.com/stretchr/testify/assert"
)

// #nosec G101 -- fake credentials, the test checks they are removed
const gemfileLockBefore = `GIT
  remote: git@git.acme.internal:platform/auth-gem.git
  revision: 9f3c1ab
  branch: feature/PROJ-123-new-login
  specs:
    acme-auth (2.1.0)
      rack (>= 2.0)
    acme-auth-extra (2.1.0)

PATH
  remote: ../shared-utils
  specs:
    acme-shared-utils (1.0.0)

GEM
  remote: https://ci-user:s3cret@gems.acme.internal/
  specs:
    acme-logger (1.4.2)
    rack (3.0.8)

GEM
  remote: https://rubygems.org/
  specs:
    nokogiri (1.16.2-x86_64-linux)
    rails (7.1.3)
      rack (>= 2.2.4)

PLATFORMS
  x86_64-linux

DEPENDENCIES
  acme-auth!
  acme-logger
  acme-shared-utils!
  rails

BUNDLED WITH
   2.5.6
`

const gemfileLockAfter = `GIT
  specs:
    acme-auth (2.1.0)

GEM
  specs:
    acme-logger (1.4.2)
    rack (3.0.8)

GEM
  specs:
    nokogiri (1.16.2-x86_64-linux)
    rails (7.1.3)
`

func TestSanitizeGemfileLock(t *testing.T) {
	sanitized, err := sanitizeGemfileLock([]byte(gemfileLockBefore), nil)

	assert.NoError(t, err)
	assert.Equal(t, gemfileLockAfter, string(sanitized))
}

func TestSanitizeGemfileLockCRLF(t *testing.T) {
	sanitized, err := sanitizeGemfileLock([]byte(strings.ReplaceAll(gemfileLockBefore, "\n", "\r\n")), nil)

	assert.NoError(t, err)
	assert.Equal(t, gemfileLockAfter, string(sanitized))
}

func TestSanitizeGemfileLockInvalid(t *testing.T) {
	for _, content := range []string{
		"",
		"  specs:\n    rack (3.0.8)\n",
		"gem 'rails'\n",
		"GEM\n  specs:\n    rack 3.0.8\n",
	} {
		_, err := sanitizeGemfileLock([]byte(content), nil)
		assert.ErrorIs(t, err, util.ErrInvalidFormat, content)
	}
}
