package npm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/debricked/cli/internal/sanitize/pm/testdata"
	"github.com/debricked/cli/internal/sanitize/pm/util"
)

func TestSanitizePackageJSON(t *testing.T) {
	sanitized, err := sanitizePackageJSON([]byte(testdata.PackageJSONBefore), newTestContext("package.json"))

	assert.NoError(t, err)
	assert.Equal(t, testdata.PackageJSONAfter, string(sanitized))
}

func TestSanitizePackageJSONInvalid(t *testing.T) {
	for _, content := range []string{"{not json", "[]", `{"a": 1} trailing`, ""} {
		_, err := sanitizePackageJSON([]byte(content), newTestContext("package.json"))
		assert.ErrorIs(t, err, util.ErrInvalidFormat, content)
	}
}

func newTestContext(path string) *util.FileContext {
	context, _ := testdata.NewContext(path)

	return context
}
