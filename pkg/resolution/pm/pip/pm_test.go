package pip

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
	assert.Equal(t, `requirements.*\.txt$`, manifest)
	_, err := regexp.Compile(manifest)
	assert.NoError(t, err)

	cases := map[string]bool{
		"requirements.txt":                         true,
		"requirements.dev.txt":                     true,
		"requirements.dev.test.txt":                true,
		"requirements-dev.test.txt":                true,
		"requirements-dev-test.txt":                true,
		"requirements-test.txt":                    true,
		"/dir/requirements.txt":                    true,
		"requirements-test-txt":                    false,
		"requirements-test.txt.dev":                false,
		"requirements-test.txt.pip.debricked.lock": false,
		"requirements.txt.pip.debricked.lock":      false,
	}
	for file, isMatch := range cases {
		t.Run(file, func(t *testing.T) {
			matched, _ := regexp.MatchString(manifest, file)
			assert.Equal(t, isMatch, matched)
		})
	}
}
