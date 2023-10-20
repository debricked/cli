package java

import (
	"fmt"
	"path/filepath"
	"testing"

	ctxTestdata "github.com/debricked/cli/internal/callgraph/cgexec/testdata"
	"github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/callgraph/finder/testdata"
	javaTestdata "github.com/debricked/cli/internal/callgraph/language/java11/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewStrategy(t *testing.T) {
	s := NewStrategy(nil, nil, nil, nil, nil, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 0)

	s = NewStrategy(nil, []string{}, []string{}, []string{}, nil, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 0)

	s = NewStrategy(nil, []string{"file"}, []string{"file"}, []string{}, nil, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 1)

	s = NewStrategy(nil, []string{"file-1", "file-2"}, []string{"file-1", "file-2"}, []string{}, nil, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 2)

	conf := config.NewConfig("java", []string{"arg1"}, map[string]string{"kwarg": "val"}, true, "maven")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1"}
	finder.FindMavenRootsNames = testFiles
	ctx, _ := ctxTestdata.NewContextMock()
	s = NewStrategy(conf, testFiles, testFiles, []string{}, finder, ctx)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 1)
	assert.Equal(t, s.config, conf)
}

func TestInvokeNoFiles(t *testing.T) {
	s := NewStrategy(nil, []string{}, []string{}, []string{}, nil, nil)
	jobs, _ := s.Invoke()
	assert.Empty(t, jobs)
}

func TestInvokeOneFile(t *testing.T) {
	conf := config.NewConfig("java", []string{"arg1"}, map[string]string{"kwarg": "val"}, true, "maven")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1"}
	finder.FindMavenRootsNames = testFiles
	ctx, _ := ctxTestdata.NewContextMock()
	s := NewStrategy(conf, testFiles, testFiles, []string{}, finder, ctx)
	jobs, _ := s.Invoke()
	assert.Len(t, jobs, 0)
}

func TestInvokeManyFiles(t *testing.T) {
	conf := config.NewConfig("java", []string{"arg1"}, map[string]string{"kwarg": "val"}, true, "maven")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1", "file-2"}
	finder.FindMavenRootsNames = testFiles
	ctx, _ := ctxTestdata.NewContextMock()
	s := NewStrategy(conf, testFiles, testFiles, []string{}, finder, ctx)
	jobs, _ := s.Invoke()
	assert.Len(t, jobs, 0)
}

func TestInvokeManyFilesWCorrectFilters(t *testing.T) {
	conf := config.NewConfig("java", []string{"arg1"}, map[string]string{"kwarg": "val"}, false, "maven")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1", "file-2", "file-3"}
	finder.FindMavenRootsNames = []string{"file-3/pom.xml"}
	finder.FindJavaClassDirsNames = []string{"file-3/test.class"}
	ctx, _ := ctxTestdata.NewContextMock()
	s := NewStrategy(conf, testFiles, testFiles, []string{"test"}, finder, ctx)
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
	conf := config.NewConfig("java", []string{"arg1"}, map[string]string{"kwarg": "val"}, false, "maven")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1", "file-2", "file-3"}
	finder.FindMavenRootsNames = []string{"file-3/pom.xml"}
	finder.FindJavaClassDirsNames = []string{"file-3/test.class"}
	ctx, _ := ctxTestdata.NewContextMock()
	s := NewStrategy(conf, testFiles, testFiles, []string{"test"}, finder, ctx)
	factoryMock := javaTestdata.NewEchoCmdFactory()
	factoryMock.BuildMavenErr = fmt.Errorf("build-error")
	s.cmdFactory = factoryMock
	err := buildProjects(s, []string{"file-3/pom.xml"})

	assert.NotNil(t, err)
}
