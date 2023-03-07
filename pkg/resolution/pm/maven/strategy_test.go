package maven

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockPomX struct{}

func (p MockPomX) GetRootPomFiles(files []string) []string {
	return files
}

func (p MockPomX) ParsePomModules(path string) ([]string, error) {
	return []string{}, nil
}

func TestNewStrategy(t *testing.T) {
	s := NewStrategy(nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 0)

	s = NewStrategy([]string{})
	assert.NotNil(t, s)
	assert.Len(t, s.files, 0)

	s = NewStrategy([]string{"file"})
	assert.NotNil(t, s)
	assert.Len(t, s.files, 1)

	s = NewStrategy([]string{"file-1", "file-2"})
	assert.NotNil(t, s)
	assert.Len(t, s.files, 2)
}

func TestInvokeNoFiles(t *testing.T) {
	s := NewStrategy([]string{})
	jobs := s.Invoke()
	assert.Empty(t, jobs)
}

func TestInvokeOneFile(t *testing.T) {
	s := NewStrategy([]string{"file"})

	s.pomX = MockPomX{}

	jobs := s.Invoke()
	assert.Len(t, jobs, 1)
}

func TestInvokeManyFiles(t *testing.T) {
	s := NewStrategy([]string{"file-1", "file-2"})

	s.pomX = MockPomX{}

	jobs := s.Invoke()
	assert.Len(t, jobs, 2)
}
