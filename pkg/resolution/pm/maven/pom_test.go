package maven

// test for pkg/resolution/pm/maven/pom.go

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePomModules(t *testing.T) {
	p := PomX{}
	modules, err := p.ParsePomModules("testdata/pom.xml")
	assert.Nil(t, err)
	assert.Len(t, modules, 5)

	modules, err = p.ParsePomModules("testdata/notAPom.xml")
	assert.NotNil(t, err)
	assert.Len(t, modules, 0)

}

func TestGetRootPomFiles(t *testing.T) {
	p := PomX{}
	files := p.GetRootPomFiles([]string{"testdata/pom.xml", "testdata/notAPom.xml"})
	assert.Len(t, files, 2)
	assert.Equal(t, "testdata/pom.xml", files[0])

	files = p.GetRootPomFiles([]string{"testdata/pom.xml", "testdata/guava/pom.xml"})
	assert.Len(t, files, 1)
	assert.Equal(t, "testdata/pom.xml", files[0])
}
