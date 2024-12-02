package java

import (
	"fmt"
	"path/filepath"
	"testing"

	ctxTestdata "github.com/debricked/cli/internal/callgraph/cgexec/testdata"
	"github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/callgraph/finder/testdata"
	javaTestdata "github.com/debricked/cli/internal/callgraph/language/java/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewStrategy(t *testing.T) {
	s := NewStrategy(nil, nil, nil, nil, nil, nil)
	assert.NotNil(t, s)

	s = NewStrategy(nil, []string{}, []string{}, []string{}, nil, nil)
	assert.NotNil(t, s)

	s = NewStrategy(nil, []string{"file"}, []string{}, []string{}, nil, nil)
	assert.NotNil(t, s)

	s = NewStrategy(nil, []string{"file-1", "file-2"}, []string{}, []string{}, nil, nil)
	assert.NotNil(t, s)

	conf := config.NewConfig("java", []string{"arg1"}, map[string]string{"kwarg": "val"}, true, "maven", "v2.0.0")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1"}
	finder.FindRootsNames = testFiles
	ctx, _ := ctxTestdata.NewContextMock()
	s = NewStrategy(conf, testFiles, []string{}, []string{}, finder, ctx)
	assert.NotNil(t, s)
	assert.Equal(t, s.config, conf)
}

func TestInvokeNoFiles(t *testing.T) {
	s := NewStrategy(nil, []string{}, []string{}, []string{}, nil, nil)
	jobs, _ := s.Invoke()
	assert.Empty(t, jobs)
}

func TestInvokeOneFile(t *testing.T) {
	conf := config.NewConfig("java", []string{"arg1"}, map[string]string{"kwarg": "val"}, true, "maven", "v2.0.0")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1"}
	finder.FindRootsNames = testFiles
	ctx, _ := ctxTestdata.NewContextMock()
	s := NewStrategy(conf, testFiles, []string{}, []string{}, finder, ctx)
	jobs, _ := s.Invoke()
	assert.Len(t, jobs, 0)
}

func TestInvokeManyFiles(t *testing.T) {
	conf := config.NewConfig("java", []string{"arg1"}, map[string]string{"kwarg": "val"}, true, "maven", "v2.0.0")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1", "file-2"}
	finder.FindRootsNames = testFiles
	ctx, _ := ctxTestdata.NewContextMock()
	s := NewStrategy(conf, testFiles, []string{}, []string{}, finder, ctx)
	jobs, _ := s.Invoke()
	assert.Len(t, jobs, 0)
}

func TestInvokeManyFilesWCorrectFilters(t *testing.T) {
	conf := config.NewConfig("java", []string{"arg1"}, map[string]string{"kwarg": "val"}, false, "maven", "v2.0.0")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1", "file-2", "file-3"}
	finder.FindRootsNames = []string{"file-3/pom.xml"}
	finder.FindDependencyDirsNames = []string{"file-3/test.class"}
	ctx, _ := ctxTestdata.NewContextMock()
	s := NewStrategy(conf, testFiles, []string{"test"}, []string{}, finder, ctx)
	jobs, _ := s.Invoke()
	assert.Len(t, jobs, 1)
	for _, job := range jobs {
		file, _ := filepath.Abs("file-3/test.class")
		dir, _ := filepath.Abs("file-3/")
		assert.Equal(t, job.GetFiles(), []string{file})
		assert.Equal(t, job.GetDir(), dir)

	}
}

func TestBuildProjectsError(t *testing.T) {
	conf := config.NewConfig("java", []string{"arg1"}, map[string]string{"kwarg": "val"}, false, "maven", "v2.0.0")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1", "file-2", "file-3"}
	finder.FindRootsNames = []string{"file-3/pom.xml"}
	finder.FindDependencyDirsNames = []string{"file-3/test.class"}
	ctx, _ := ctxTestdata.NewContextMock()
	s := NewStrategy(conf, testFiles, []string{"test"}, []string{}, finder, ctx)
	factoryMock := javaTestdata.NewEchoCmdFactory()
	factoryMock.BuildMavenErr = fmt.Errorf("build-error")
	s.cmdFactory = factoryMock
	err := buildProjects(s, []string{"file-3/pom.xml"})

	assert.NotNil(t, err)
}
