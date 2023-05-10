package finder

import (
	"path/filepath"
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

func TestGetRootPomFiles(t *testing.T) {
	pomParent := filepath.Join("testdata", "pom.xml")
	pomFail := filepath.Join("testdata", "notAPom.xml")
	pomChild := filepath.Join("testdata", "guava", "pom.xml")

	p := PomService{}
	files := p.GetRootPomFiles([]string{pomParent, pomFail})
	assert.Len(t, files, 1)

	files = p.GetRootPomFiles([]string{pomParent, pomChild})
	assert.Len(t, files, 1)
	assert.Equal(t, pomParent, files[0])

	files = p.GetRootPomFiles([]string{pomFail})
	assert.Len(t, files, 0)
}
