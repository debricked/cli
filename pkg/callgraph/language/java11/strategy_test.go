package java

import (
	"testing"

	"github.com/debricked/cli/pkg/callgraph/config"
	"github.com/debricked/cli/pkg/io/finder/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewStrategy(t *testing.T) {
	s := NewStrategy(nil, nil, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 0)

	s = NewStrategy(nil, []string{}, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 0)

	s = NewStrategy(nil, []string{"file"}, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 1)

	s = NewStrategy(nil, []string{"file-1", "file-2"}, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 2)

	conf := config.NewConfig("java", []string{"arg1"}, map[string]string{"kwarg": "val"})
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1"}
	finder.FindMavenRootsNames = testFiles
	s = NewStrategy(conf, testFiles, finder)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 2)
	assert.Equal(t, s.config, conf)
}

func TestInvokeNoFiles(t *testing.T) {
	s := NewStrategy(nil, []string{}, nil)
	jobs, _ := s.Invoke()
	assert.Empty(t, jobs)
}

func TestInvokeOneFile(t *testing.T) {
	conf := config.NewConfig("java", []string{"arg1"}, map[string]string{"kwarg": "val"})
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1"}
	finder.FindMavenRootsNames = testFiles
	s := NewStrategy(conf, testFiles, finder)
	jobs, _ := s.Invoke()
	assert.Len(t, jobs, 1)
}

func TestInvokeManyFiles(t *testing.T) {
	conf := config.NewConfig("java", []string{"arg1"}, map[string]string{"kwarg": "val"})
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1", "file-2"}
	finder.FindMavenRootsNames = testFiles
	s := NewStrategy(conf, testFiles, finder)
	jobs, _ := s.Invoke()
	assert.Len(t, jobs, 2)
}

func TestInvokeManyFilesWFilters(t *testing.T) {
	conf := config.NewConfig("java", []string{"arg1"}, map[string]string{"kwarg": "val"})
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1", "file-2"}
	finder.FindMavenRootsNames = []string{"file-3"}
	s := NewStrategy(conf, testFiles, finder)
	jobs, _ := s.Invoke()
	assert.Equal(t, jobs, []string{"file-3"})
	assert.Len(t, jobs, 1)
}
