package pm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPms(t *testing.T) {
	pms := Pms()
	pmNames := []string{
		"mvn",
		"go",
		"gradle",
	}

	for _, pmName := range pmNames {
		t.Run(pmName, func(t *testing.T) {
			contains := false
			for _, pm := range pms {
				contains = contains || pm.Name() == pmName
			}
			assert.Truef(t, contains, "failed to assert that %s was returned in Pms()", pmName)
		})
	}
}
