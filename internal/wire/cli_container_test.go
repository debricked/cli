package wire

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWire(t *testing.T) {
	cliContainer = &CliContainer{}
	defer resetContainer()

	err := cliContainer.wire()
	assert.NoError(t, err)
	assertCliContainer(t, cliContainer)
}

func TestGetCliContainer(t *testing.T) {
	assert.Nil(t, cliContainer)
	testGetCliContainer(t)
}

func testGetCliContainer(t *testing.T) {
	container := GetCliContainer()
	assert.NotNil(t, container)
	assert.NotNil(t, cliContainer)
	assertCliContainer(t, cliContainer)
}

func resetContainer() {
	cliContainer = nil
}

func assertCliContainer(t *testing.T, cc *CliContainer) {
	assert.NotNil(t, cc.DebClient())
	assert.NotNil(t, cc.Finder())
	assert.NotNil(t, cc.Scanner())
	assert.NotNil(t, cc.Resolver())
	assert.NotNil(t, cc.CallgraphGenerator())
	assert.NotNil(t, cc.LicenseReporter())
	assert.NotNil(t, cc.VulnerabilityReporter())
	assert.NotNil(t, cc.Fingerprinter())
	assert.NotNil(t, cc.Authenticator())
	assert.NotNil(t, cc.SBOMReporter())
}
