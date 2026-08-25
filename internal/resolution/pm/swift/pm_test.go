package swift

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestName(t *testing.T) {
	pm := NewPm()
	assert.Equal(t, Name, pm.Name())
}

func TestManifests(t *testing.T) {
	pm := NewPm()
	manifests := pm.Manifests()
	assert.Len(t, manifests, 1)
	assert.Equal(t, `Package\.swift$`, manifests[0])
}
