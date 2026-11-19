package nuget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/debricked/cli/internal/sanitize/pm/testdata"
	"github.com/debricked/cli/internal/sanitize/pm/util"
)

const packagesLockBefore = `{
  "version": 2,
  "dependencies": {
    "net8.0": {
      "Acme.Logger": {
        "type": "Direct",
        "requested": "[1.4.2, )",
        "resolved": "1.4.2",
        "contentHash": "AAAA",
        "dependencies": {
          "Newtonsoft.Json": "13.0.1"
        }
      },
      "acme.payments.core": {
        "type": "Project",
        "dependencies": {
          "System.Buffers": "[4.5.1, )"
        }
      }
    },
    "netstandard2.0": {
      "Acme.Logger": {
        "type": "Direct",
        "requested": "[1.4.2, )",
        "resolved": "1.4.2",
        "contentHash": "AAAA",
        "dependencies": {
          "Acme.Payments.Core": "1.0.0",
          "Missing.From.Group": "2.0.0"
        }
      },
      "acme.payments.core": {
        "type": "project",
        "dependencies": {
          "System.Buffers": "[4.5.1, )"
        }
      },
      "System.Buffers": {
        "type": "Transitive",
        "resolved": "4.5.1",
        "contentHash": "CCCC"
      }
    }
  }
}
`

const packagesLockAfter = `{
    "dependencies": {
        "net8.0": {
            "Acme.Logger": {
                "type": "Direct",
                "requested": "[1.4.2, )",
                "resolved": "1.4.2",
                "contentHash": "AAAA",
                "dependencies": {
                    "Newtonsoft.Json": "13.0.1"
                }
            },
              "debricked-redacted-project-1": {
                "type": "Project",
                "dependencies": {
                    "System.Buffers": "[4.5.1, )"
                }
            }
        },
        "netstandard2.0": {
            "Acme.Logger": {
                "type": "Direct",
                "requested": "[1.4.2, )",
                "resolved": "1.4.2",
                "contentHash": "AAAA",
                "dependencies": {
                  "debricked-redacted-project-1": "1.0.0",
                    "Missing.From.Group": "2.0.0"
                }
            },
              "debricked-redacted-project-1": {
                "type": "project",
                "dependencies": {
                    "System.Buffers": "[4.5.1, )"
                }
            },
            "System.Buffers": {
                "type": "Transitive",
                "resolved": "4.5.1",
                "contentHash": "CCCC"
            }
        }
    }
}
`

func TestSanitizePackagesLockJSON(t *testing.T) {
	sanitized, err := sanitizePackagesLockJSON([]byte(packagesLockBefore), newTestContext("packages.lock.json"))

	assert.NoError(t, err)
	assert.JSONEq(t, packagesLockAfter, string(sanitized))
}

func TestSanitizePackagesLockJSONPreservesOneKeyPerLine(t *testing.T) {
	sanitized, err := sanitizePackagesLockJSON([]byte(packagesLockBefore), newTestContext("packages.lock.json"))

	assert.NoError(t, err)
	assert.Contains(t, string(sanitized), "\n                \"contentHash\": \"AAAA\",\n")
}

func TestSanitizePackagesLockJSONRenamesProjectEntriesCaseInsensitive(t *testing.T) {
	context, placeholders := testdata.NewContext("packages.lock.json")

	sanitized, err := sanitizePackagesLockJSON([]byte(packagesLockBefore), context)

	assert.NoError(t, err)
	assert.Contains(t, string(sanitized), `"debricked-redacted-project-1"`)
	assert.Len(t, placeholders.Entries(), 1)
	assert.Equal(t, util.MapEntry{
		Placeholder: "debricked-redacted-project-1",
		Name:        "acme.payments.core",
		Path:        "packages.lock.json",
	}, placeholders.Entries()[0])
}

func TestSanitizePackagesLockJSONKeepsOtherPackagesNamedLikeAPlaceholder(t *testing.T) {
	before := `{
  "dependencies": {
    "net8.0": {
      "redacted-project-1": {
        "type": "Direct",
        "resolved": "1.0.0",
        "contentHash": "AAAA"
      },
      "Acme.Payments.Core": {
        "type": "Project",
        "dependencies": {
          "redacted-project-1": "1.0.0"
        }
      }
    }
  }
}
`

	sanitized, err := sanitizePackagesLockJSON([]byte(before), newTestContext("packages.lock.json"))

	assert.NoError(t, err)
	assert.Contains(t, string(sanitized), `"redacted-project-1": {
                "type": "Direct"`)
	assert.Contains(t, string(sanitized), `"debricked-redacted-project-1": {
                "type": "Project"`)
	assert.NotContains(t, string(sanitized), "Acme.Payments.Core")
}

func TestSanitizePackagesLockJSONInvalid(t *testing.T) {
	invalid := []string{
		`{not-json`,
		`[]`,
		`{"dependencies": []}`,
		`{"dependencies": {"net8.0": []}}`,
		`{"dependencies": {"net8.0": {"pkg": []}}}`,
		`{"dependencies": {"net8.0": {"pkg": {"resolved": "1.0.0"}}}}`,
		`{"dependencies": {"net8.0": {"pkg": {"type": 1}}}}`,
		`{"dependencies": {"net8.0": {"pkg": {"type": "Direct", "resolved": {"nested": "x"}}}}}`,
		`{"dependencies": {"net8.0": {"pkg": {"type": "Direct", "dependencies": []}}}}`,
		`{"dependencies": {"net8.0": {"pkg": {"type": "Direct", "dependencies": {"a": 1}}}}}`,
	}
	for _, content := range invalid {
		_, err := sanitizePackagesLockJSON([]byte(content), newTestContext("packages.lock.json"))
		assert.ErrorIs(t, err, util.ErrInvalidFormat, content)
	}
}
