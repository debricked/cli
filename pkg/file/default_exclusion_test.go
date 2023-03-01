package file

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultExclusions(t *testing.T) {
	separator := string(os.PathSeparator)
	for _, ex := range DefaultExclusions() {
		exParts := strings.Split(ex, separator)
		assert.Greaterf(t, len(exParts), 0, "failed to assert that %s used correct separator. Proper separator %s", ex, separator)
	}
}
