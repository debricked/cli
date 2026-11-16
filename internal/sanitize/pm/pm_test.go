package pm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPms(t *testing.T) {
	names := []string{}
	for _, pm := range Pms() {
		names = append(names, pm.Name())
	}

	assert.Equal(t, []string{"npm", "yarn", "pnpm", "nuget"}, names)
}

func TestEachFileNameIsHandledByOnePm(t *testing.T) {
	owners := map[string]string{}
	for _, pm := range Pms() {
		for _, files := range []map[string]any{toAny(pm.LockFiles()), toAny(pm.Manifests())} {
			for name := range files {
				owner, claimed := owners[name]
				assert.False(t, claimed, "%s is handled by %s and %s", name, owner, pm.Name())
				owners[name] = pm.Name()
			}
		}
	}
}

func toAny[V any](files map[string]V) map[string]any {
	converted := map[string]any{}
	for name, value := range files {
		converted[name] = value
	}

	return converted
}
