package cocoapods

import (
	"testing"

	"github.com/debricked/cli/internal/sanitize/pm/util"
	"github.com/stretchr/testify/assert"
)

const podfileLockBefore = `PODS:
  - AcmeAuth (2.1.0):
    - Alamofire (~> 5.0)
  - AcmeLogger (1.4.2)
  - Alamofire (5.8.1)
  - Firebase/Core (10.20.0):
    - Firebase/CoreOnly
  - Firebase/CoreOnly (10.20.0)
  - "GoogleUtilities/Environment (7.12.0)":
    - PromisesObjC (< 3.0, >= 1.2)

DEPENDENCIES:
  - AcmeAuth (from ` + "`git@git.acme.internal:ios/auth.git`" + `, branch ` + "`feature/PROJ-123`" + `)
  - AcmeLogger (from ` + "`../AcmeLogger`" + `)
  - Alamofire
  - Firebase/Core

SPEC REPOS:
  "https://cdn.cocoapods.org/":
    - Alamofire
  "https://git.acme.internal/ios/specs.git":
    - Firebase

EXTERNAL SOURCES:
  AcmeAuth:
    :branch: feature/PROJ-123
    :git: "git@git.acme.internal:ios/auth.git"
  AcmeLogger:
    :path: "../AcmeLogger"

CHECKOUT OPTIONS:
  AcmeAuth:
    :commit: 9f3c1ab
    :git: "git@git.acme.internal:ios/auth.git"

SPEC CHECKSUMS:
  AcmeAuth: aaaa
  Alamofire: bbbb

PODFILE CHECKSUM: cccc

COCOAPODS: 1.15.2
`

const podfileLockAfter = `PODS:
  - AcmeAuth (2.1.0)
  - AcmeLogger (1.4.2)
  - Alamofire (5.8.1)
  - Firebase/Core (10.20.0)
  - Firebase/CoreOnly (10.20.0)
  - "GoogleUtilities/Environment (7.12.0)"

EXTERNAL SOURCES:
  AcmeAuth:
  AcmeLogger:

`

func TestSanitizePodfileLock(t *testing.T) {
	sanitized, err := sanitizePodfileLock([]byte(podfileLockBefore), nil)

	assert.NoError(t, err)
	assert.Equal(t, podfileLockAfter, string(sanitized))
	assert.NotContains(t, string(sanitized), "acme.internal")
}

func TestSanitizePodfileLockWithoutExternalSources(t *testing.T) {
	sanitized, err := sanitizePodfileLock([]byte("PODS:\n  - Alamofire (5.8.1)\n\nCOCOAPODS: 1.15.2\n"), nil)

	assert.NoError(t, err)
	assert.Equal(t, "PODS:\n  - Alamofire (5.8.1)\n\n", string(sanitized))
}

func TestSanitizePodfileLockInvalid(t *testing.T) {
	for _, content := range []string{
		"",
		"COCOAPODS: 1.15.2\n",
		"  - Alamofire (5.8.1)\n",
		"pods:\n  - Alamofire (5.8.1)\n",
		"PODS:\n  - Alamofire\n",
		"PODS:\n  - Alamofire (5.8.1)\nEXTERNAL SOURCES:\n  AcmeAuth: {git: x}\n",
	} {
		_, err := sanitizePodfileLock([]byte(content), nil)
		assert.ErrorIs(t, err, util.ErrInvalidFormat, content)
	}
}
