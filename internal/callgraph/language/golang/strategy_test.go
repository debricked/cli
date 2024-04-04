package golang

import (
	"testing"

	ctxTestdata "github.com/debricked/cli/internal/callgraph/cgexec/testdata"
	"github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/callgraph/finder/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewStrategy(t *testing.T) {
	s := NewStrategy(nil, nil, nil, nil, nil)
	assert.NotNil(t, s)

	s = NewStrategy(nil, []string{}, []string{}, nil, nil)
	assert.NotNil(t, s)

	s = NewStrategy(nil, []string{"file"}, []string{}, nil, nil)
	assert.NotNil(t, s)

	s = NewStrategy(nil, []string{"file-1", "file-2"}, []string{}, nil, nil)
	assert.NotNil(t, s)

	conf := config.NewConfig("golang", []string{"arg1"}, map[string]string{"kwarg": "val"}, true, "go")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1"}
	finder.FindRootsNames = testFiles
	ctx, _ := ctxTestdata.NewContextMock()
	s = NewStrategy(conf, []string{"."}, []string{}, finder, ctx)
	assert.NotNil(t, s)
	assert.Equal(t, s.config, conf)
}

func TestInvokeNoFiles(t *testing.T) {
	s := NewStrategy(nil, []string{}, []string{}, nil, nil)
	jobs, _ := s.Invoke()
	assert.Empty(t, jobs)
}

func TestInvokeOneFile(t *testing.T) {
	conf := config.NewConfig("golang", []string{"arg1"}, map[string]string{"kwarg": "val"}, true, "go")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1"}
	finder.FindRootsNames = testFiles
	ctx, _ := ctxTestdata.NewContextMock()
	s := NewStrategy(conf, []string{"."}, []string{}, finder, ctx)
	jobs, err := s.Invoke()
	assert.NoError(t, err)
	assert.Len(t, jobs, 1)
}

func TestInvokeManyFiles(t *testing.T) {
	conf := config.NewConfig("golang", []string{"arg1"}, map[string]string{"kwarg": "val"}, true, "go")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1", "file-2"}
	finder.FindRootsNames = testFiles
	ctx, _ := ctxTestdata.NewContextMock()
	s := NewStrategy(conf, []string{"."}, []string{}, finder, ctx)
	jobs, _ := s.Invoke()
	assert.Len(t, jobs, 2)
}

func TestInvokeWithErrors(t *testing.T) {
	conf := config.NewConfig("golang", []string{"arg1"}, map[string]string{"kwarg": "val"}, true, "go")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{"file-1", "file-2"}
	finder.FindRootsNames = testFiles
	finder.FindRootsErr = assert.AnError
	ctx, _ := ctxTestdata.NewContextMock()
	s := NewStrategy(conf, []string{"."}, []string{}, finder, ctx)
	jobs, err := s.Invoke()
	assert.Error(t, err)
	assert.Empty(t, jobs)

	finder.FindRootsErr = nil
	finder.FindFilesErr = assert.AnError
	s = NewStrategy(conf, []string{"."}, []string{}, finder, ctx)
	jobs, err = s.Invoke()
	assert.Error(t, err)
	assert.Empty(t, jobs)
}

func TestInvokeNoRoots(t *testing.T) {
	conf := config.NewConfig("golang", []string{"arg1"}, map[string]string{"kwarg": "val"}, true, "go")
	finder := testdata.NewEmptyFinderMock()
	testFiles := []string{}
	finder.FindRootsNames = testFiles
	ctx, _ := ctxTestdata.NewContextMock()
	s := NewStrategy(conf, []string{"."}, []string{}, finder, ctx)
	jobs, err := s.Invoke()
	assert.NoError(t, err)
	assert.Empty(t, jobs)
}
