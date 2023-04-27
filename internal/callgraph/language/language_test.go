package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLanguages(t *testing.T) {
	langs := Languages()
	langNames := []string{
		"java",
	}

	for _, langName := range langNames {
		t.Run(langName, func(t *testing.T) {
			contains := false
			for _, pm := range langs {
				contains = contains || pm.Name() == langName
			}
			assert.Truef(t, contains, "failed to assert that %s was returned in Languages()", langName)
		})
	}
}
