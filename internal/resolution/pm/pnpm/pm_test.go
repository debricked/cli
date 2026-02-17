package pnpm

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

	// First manifest: package.json
	manifestPkg := manifests[0]
	assert.Equal(t, `package\.json$`, manifestPkg)
	_, err := regexp.Compile(manifestPkg)
	assert.NoError(t, err)

	casesPkg := map[string]bool{
		"package.json":      true,
		"pnpm-lock.yaml":    false,
		"pnpm-lock.yml":     false,
		"package-lock.json": false,
	}
	for file, isMatch := range casesPkg {
		t.Run("pkg-"+file, func(t *testing.T) {
			matched, _ := regexp.MatchString(manifestPkg, file)
			assert.Equal(t, isMatch, matched)
		})
	}

	// Second manifest: pnpm-lock.yaml
	manifestLock := manifests[1]
	assert.Equal(t, `pnpm-lock\.yaml$`, manifestLock)
	_, err = regexp.Compile(manifestLock)
	assert.NoError(t, err)

	casesLock := map[string]bool{
		"pnpm-lock.yaml": true,
		"pnpm-lock.yml":  false,
		"package.json":   false,
	}
	for file, isMatch := range casesLock {
		t.Run("lock-"+file, func(t *testing.T) {
			matched, _ := regexp.MatchString(manifestLock, file)
			assert.Equal(t, isMatch, matched)
		})
	}
}
