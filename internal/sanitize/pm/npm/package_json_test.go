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

// A yarn.lock dependency on the same patch gets the same value, see the yarn tests.
func TestSanitizePackageJSONKeepsThePatchedDescriptor(t *testing.T) {
	content := `{"dependencies": {"interpret": "patch:interpret@npm%3A3.1.1#~/.yarn/patches/interpret-npm-3.1.1-715bac2bd7.patch"}}`

	sanitized, err := sanitizePackageJSON([]byte(content), newTestContext("package.json"))

	assert.NoError(t, err)
	assert.Contains(t, string(sanitized), `"interpret": "patch:interpret@npm%3A3.1.1#redacted.patch"`)
	assert.NotContains(t, string(sanitized), ".yarn/patches")
}
