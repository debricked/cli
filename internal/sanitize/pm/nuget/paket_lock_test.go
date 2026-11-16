package nuget

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/debricked/cli/internal/sanitize/pm/testdata"
	"github.com/debricked/cli/internal/sanitize/pm/util"
)

func TestSanitizePaketLock(t *testing.T) {
	sanitized, err := sanitizePaketLock([]byte(testdata.PaketBefore), newTestContext("paket.lock"))

	assert.NoError(t, err)
	assert.Equal(t, testdata.PaketAfter, string(sanitized))
}

func TestSanitizePaketLockInvalid(t *testing.T) {
	_, err := sanitizePaketLock([]byte("this is not a lock file\n"), newTestContext("paket.lock"))

	assert.ErrorIs(t, err, util.ErrInvalidFormat)
}

func TestSanitizePaketLockAvoidsPlaceholderCollisions(t *testing.T) {
	content := `NUGET
    redacted-file-1 (1.0.0)
GITHUB
    src/deploy/Release.fs (9f3c1ab)
`

	sanitized, err := sanitizePaketLock([]byte(content), newTestContext("paket.lock"))

	assert.NoError(t, err)
	assert.Equal(t, "    redacted-file-1 (1.0.0)\n    redacted-file-2 (9f3c1ab)\n", string(sanitized))
}

func newTestContext(path string) *util.FileContext {
	context, _ := testdata.NewContext(path)

	return context
}
