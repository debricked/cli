package dart

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLanguage(t *testing.T) {
	lang := NewLanguage()
	assert.Equal(t, Name, lang.name)
	assert.Equal(t, StandardVersion, lang.version)
}

func TestName(t *testing.T) {
	lang := NewLanguage()
	assert.Equal(t, Name, lang.Name())
}

func TestVersion(t *testing.T) {
	lang := NewLanguage()
	assert.Equal(t, StandardVersion, lang.Version())
}
