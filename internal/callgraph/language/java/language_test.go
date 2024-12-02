package java

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLanguage(t *testing.T) {
	pm := NewLanguage()
	assert.Equal(t, Name, pm.name)
	assert.Equal(t, StandardVersion, pm.version)
}

func TestName(t *testing.T) {
	pm := NewLanguage()
	assert.Equal(t, Name, pm.Name())
}

func TestVersion(t *testing.T) {
	pm := NewLanguage()
	assert.Equal(t, StandardVersion, pm.Version())
}
