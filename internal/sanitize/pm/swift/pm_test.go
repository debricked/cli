package swift

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPm(t *testing.T) {
	pm := NewPm()

	assert.Equal(t, Name, pm.Name())
	assert.Len(t, pm.LockFiles(), 1)
	assert.Contains(t, pm.LockFiles(), "Package.resolved")
	assert.Empty(t, pm.Manifests())
}
