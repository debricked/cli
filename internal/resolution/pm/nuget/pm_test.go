package nuget

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
	assert.Len(t, manifests, 2)
	manifestCs := manifests[0]
	assert.Equal(t, `\.csproj$`, manifestCs)
	_, err := regexp.Compile(manifestCs)
	assert.NoError(t, err)

	manifestPc := manifests[1]
	assert.Equal(t, `packages\.config$`, manifestPc)
	_, err = regexp.Compile(manifestPc)
	assert.NoError(t, err)

	cases := map[string]bool{
		"test.csproj":             true,
		"sample3.csproj":          true,
		".csproj":                 true,
		"test.csproj.user":        false,
		"test.csproj.nuget":       false,
		"test.csproj.nuget.props": false,
		"package.json.lock":       false,
		"packages.config":         true,
	}
	for file, isMatch := range cases {
		t.Run(file, func(t *testing.T) {

			matchedCs, _ := regexp.MatchString(manifestCs, file)
			matchedPc, _ := regexp.MatchString(manifestPc, file)
			matched := matchedCs || matchedPc

			assert.Equal(t, isMatch, matched, "file: %s", file)
		})
	}
}
