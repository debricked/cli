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

func newTestContext(path string) *util.FileContext {
	context, _ := testdata.NewContext(path)

	return context
}
