package uv

import "testing"

func TestManifests(t *testing.T) {
	pm := NewPm()
	manifests := pm.Manifests()

	if len(manifests) != 1 {
		t.Fatalf("expected 1 manifest, got %d", len(manifests))
	}

	if manifests[0] != "pyproject.toml" {
		t.Fatalf("expected pyproject.toml, got %s", manifests[0])
	}
}
