package file

import (
	"os"
	"strings"
	"testing"
)

func TestDefaultExclusions(t *testing.T) {
	separator := string(os.PathSeparator)
	for _, ex := range DefaultExclusions() {
		exParts := strings.Split(ex, separator)
		if len(exParts) == 0 {
			t.Errorf("failed to assert that %s used correct separator. Proper separator %s", ex, separator)
		}
	}
}
