package swift

import (
	"testing"

	"github.com/debricked/cli/internal/sanitize/pm/util"
	"github.com/stretchr/testify/assert"
)

const packageResolvedBefore = `{
  "originHash": "aaaa",
  "pins": [
    {
      "identity": "acme-auth",
      "kind": "remoteSourceControl",
      "location": "git@git.acme.internal:ios/acme-auth.git",
      "state": {
        "branch": "feature/PROJ-123-new-login",
        "revision": "9f3c1ab"
      }
    },
    {
      "identity": "swift-log",
      "kind": "remoteSourceControl",
      "location": "https://github.com/apple/swift-log.git",
      "state": {
        "revision": "bbbb",
        "version": "1.5.4"
      }
    }
  ],
  "version": 2
}`

const packageResolvedAfter = `{
    "pins": [
        {
            "identity": "acme-auth",
            "state": {
                "revision": "9f3c1ab"
            }
        },
        {
            "identity": "swift-log",
            "state": {
                "version": "1.5.4"
            }
        }
    ]
}
`

func TestSanitizePackageResolved(t *testing.T) {
	sanitized, err := sanitizePackageResolved([]byte(packageResolvedBefore), nil)

	assert.NoError(t, err)
	assert.Equal(t, packageResolvedAfter, string(sanitized))
}

func TestSanitizePackageResolvedWithoutIdentityKeepsOnlyTheLastPartOfTheLocation(t *testing.T) {
	before := `{"pins": [
		{"location": "git@git.acme.internal:ios/acme-auth.git", "state": {"branch": "main", "version": null, "revision": "9f3c1ab"}},
		{"location": "https://github.com/apple/swift-log/", "state": {}},
		{"identity": "no-state"}
	]}`

	sanitized, err := sanitizePackageResolved([]byte(before), nil)

	assert.NoError(t, err)
	assert.Equal(t, `{
    "pins": [
        {
            "location": "https://redacted.invalid/acme-auth.git",
            "state": {
                "revision": "9f3c1ab"
            }
        },
        {
            "location": "https://redacted.invalid/swift-log",
            "state": {}
        },
        {
            "identity": "no-state",
            "state": {}
        }
    ]
}
`, string(sanitized))
}

func TestSanitizePackageResolvedInvalid(t *testing.T) {
	for _, content := range []string{
		``,
		`[]`,
		`{"object": {"pins": []}, "version": 1}`,
		`{"pins": {}}`,
		`{"pins": ["acme-auth"]}`,
		`{"pins": [{"kind": "remoteSourceControl"}]}`,
		`{"pins": [{"location": "/"}]}`,
	} {
		_, err := sanitizePackageResolved([]byte(content), nil)
		assert.ErrorIs(t, err, util.ErrInvalidFormat, content)
	}
}
