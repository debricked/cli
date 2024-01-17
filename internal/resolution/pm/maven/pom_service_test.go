package maven

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePomModules(t *testing.T) {
	p := PomService{}
	modules, err := p.ParsePomModules("testdata/pom.xml")
	assert.Nil(t, err)
	assert.Len(t, modules, 5)
	correct := []string{"guava", "guava-bom", "guava-gwt", "guava-testlib", "guava-tests"}
	assert.Equal(t, correct, modules)

	modules, err = p.ParsePomModules("testdata/notAPom.xml")

	assert.NotNil(t, err)
	assert.Len(t, modules, 0)
}
