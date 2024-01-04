package bower

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPm(t *testing.T) {
	pm := NewPm()
	assert.Equal(t, Name, pm.name)
}

func TestName(t *testing.T) {
	pm := NewPm()
	assert.Equal(t, Name, pm.Name())
}

func TestManifests(t *testing.T) {
	pm := Pm{}
	manifests := pm.Manifests()
	assert.Len(t, manifests, 1)
	manifest := manifests[0]
	assert.Equal(t, `bower\.json$`, manifest)
	_, err := regexp.Compile(manifest)
	assert.NoError(t, err)

	cases := map[string]bool{
		"bower.json":        true,
		"package.json":      false,
		"package-lock.json": false,
		"bower.lock":        false,
	}
	for file, isMatch := range cases {
		t.Run(file, func(t *testing.T) {
			matched, _ := regexp.MatchString(manifest, file)
			assert.Equal(t, isMatch, matched)
		})
	}
}
