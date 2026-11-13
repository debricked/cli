package npm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/debricked/cli/internal/sanitize/pm/util"
)

const packageLockBefore = `{
  "name": "@acme/payments-portal",
  "version": "3.4.1",
  "lockfileVersion": 3,
  "requires": true,
  "packages": {
    "": {
      "name": "@acme/payments-portal",
      "version": "3.4.1",
      "license": "UNLICENSED",
      "dependencies": {
        "express": "^4.18.2",
        "@acme/logger": "^1.4.0",
        "@acme/auth-client": "git+ssh://git@git.acme.internal/platform/auth-client.git#v2.1.0",
        "@acme/shared-utils": "file:../shared-utils",
        "@acme/legacy": "git+ssh://git@git.acme.internal/platform/authclient.git#v1.0.0"
      },
      "devDependencies": {
        "jest": "^29.7.0"
      }
    },
    "../shared-utils": {
      "version": "1.0.0",
      "license": "UNLICENSED"
    },
    "node_modules/@acme/auth-client": {
      "version": "2.1.0",
      "resolved": "git+ssh://git@git.acme.internal/platform/auth-client.git#9f3c1ab",
      "license": "UNLICENSED"
    },
    "node_modules/@acme/logger": {
      "version": "1.4.2",
      "resolved": "https://artifactory.acme.internal/api/npm/npm-local/@acme/logger/-/logger-1.4.2.tgz",
      "integrity": "sha512-AAAA",
      "license": "UNLICENSED"
    },
    "node_modules/@acme/shared-utils": {
      "resolved": "../shared-utils",
      "link": true
    },
    "node_modules/express": {
      "version": "4.18.2",
      "resolved": "https://registry.npmjs.org/express/-/express-4.18.2.tgz",
      "integrity": "sha512-BBBB",
      "license": "MIT",
      "dependencies": {
        "qs": "6.11.0"
      },
      "engines": {
        "node": ">= 0.10.0"
      }
    },
    "node_modules/jest": {
      "version": "29.7.0",
      "resolved": "https://registry.npmjs.org/jest/-/jest-29.7.0.tgz",
      "integrity": "sha512-CCCC",
      "dev": true,
      "license": "MIT",
      "bin": {
        "jest": "bin/jest.js"
      }
    },
    "node_modules/qs": {
      "version": "6.11.0",
      "resolved": "https://registry.npmjs.org/qs/-/qs-6.11.0.tgz",
      "integrity": "sha512-DDDD",
      "license": "BSD-3-Clause"
    },
    "node_modules/@acme/legacy": {
      "version": "1.0.0",
      "resolved": "git+ssh://git@git.acme.internal/platform/authclient.git#v1.0.0",
      "license": "UNLICENSED"
    }
  }
}
`

const packageLockAfter = `{
    "lockfileVersion": 3,
    "packages": {
        "": {
            "dependencies": {
                "express": "^4.18.2",
                "@acme/logger": "^1.4.0",
                "@acme/auth-client": "git+ssh://git@redacted/redacted-repo.git#redacted",
                "@acme/shared-utils": "file:./redacted",
                "@acme/legacy": "git+ssh://git@redacted/redacted.git#redacted"
            },
            "devDependencies": {
                "jest": "^29.7.0"
            }
        },
        "./redacted-1": {
            "version": "1.0.0"
        },
        "node_modules/@acme/auth-client": {
            "version": "2.1.0"
        },
        "node_modules/@acme/logger": {
            "version": "1.4.2"
        },
        "node_modules/@acme/shared-utils": {
            "resolved": "./redacted-1",
            "link": true
        },
        "node_modules/express": {
            "version": "4.18.2",
            "dependencies": {
                "qs": "6.11.0"
            }
        },
        "node_modules/jest": {
            "version": "29.7.0",
            "dev": true
        },
        "node_modules/qs": {
            "version": "6.11.0"
        },
        "node_modules/@acme/legacy": {
            "version": "1.0.0",
            "resolved": "git+ssh://git@redacted/redacted.git#redacted"
        }
    }
}
`

func TestSanitizePackageLock(t *testing.T) {
	sanitized, err := sanitizePackageLock([]byte(packageLockBefore), newTestContext("package-lock.json"))

	assert.NoError(t, err)
	assert.Equal(t, packageLockAfter, string(sanitized))
}

func TestSanitizePackageLockVersion1(t *testing.T) {
	before := `{
  "name": "portal",
  "lockfileVersion": 1,
  "requires": true,
  "dependencies": {
    "express": {
      "version": "https://registry.npmjs.org/express/-/express-4.18.2.tgz",
      "from": "express@latest",
      "integrity": "sha512-A",
      "dev": true,
      "requires": {
        "qs": "^6.0.0",
        "acme": "git+ssh://git@git.acme.internal/p/acme.git"
      },
      "dependencies": {
        "qs": {
          "version": "6.11.0",
          "resolved": "https://artifactory.acme.internal/qs.tgz"
        }
      }
    },
    "local": {
      "version": "file:../local"
    }
  }
}
`
	after := `{
    "lockfileVersion": 1,
    "dependencies": {
        "express": {
            "version": "https://registry.npmjs.org/express/-/express-4.18.2.tgz",
            "dev": true,
            "requires": {
                "qs": "^6.0.0",
                "acme": "git+ssh://git@redacted/redacted.git#redacted"
            },
            "dependencies": {
                "qs": {
                    "version": "6.11.0"
                }
            }
        },
        "local": {
            "version": "file:./redacted"
        }
    }
}
`
	sanitized, err := sanitizePackageLock([]byte(before), newTestContext("package-lock.json"))

	assert.NoError(t, err)
	assert.Equal(t, after, string(sanitized))
}

func TestSanitizePackageLockKeepsPackageNameUnderWorkspace(t *testing.T) {
	before := `{
  "lockfileVersion": 3,
  "packages": {
    "packages/billing": {"name": "@acme/billing", "version": "1.0.0"},
    "packages/billing/node_modules/qs": {"version": "6.11.0"},
    "node_modules/@acme/billing": {"resolved": "packages/billing", "link": true},
    "node_modules/../../escape": {"version": "1.0.0"}
  }
}
`
	after := `{
    "lockfileVersion": 3,
    "packages": {
        "./redacted-1": {
            "version": "1.0.0"
        },
        "./redacted-1/node_modules/qs": {
            "version": "6.11.0"
        },
        "node_modules/@acme/billing": {
            "resolved": "./redacted-1",
            "link": true
        },
        "./redacted-2": {
            "version": "1.0.0"
        }
    }
}
`
	sanitized, err := sanitizePackageLock([]byte(before), newTestContext("package-lock.json"))

	assert.NoError(t, err)
	assert.Equal(t, after, string(sanitized))
}

func TestSanitizePackageLockInvalid(t *testing.T) {
	invalid := []string{"{not json", "[]", `{"name": "not a lock file"}`, `{"packages": []}`, `{"dependencies": "x"}`}
	for _, content := range invalid {
		_, err := sanitizePackageLock([]byte(content), newTestContext("package-lock.json"))
		assert.ErrorIs(t, err, util.ErrInvalidFormat, content)
	}
}
