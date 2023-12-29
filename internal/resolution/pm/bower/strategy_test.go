package bower

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
	jobs, _ := s.Invoke()
	assert.Empty(t, jobs)
}

func TestInvokeOneFile(t *testing.T) {
	s := NewStrategy([]string{"file"})
	jobs, _ := s.Invoke()
	assert.Len(t, jobs, 1)
}

func TestInvokeManyFiles(t *testing.T) {
	s := NewStrategy([]string{"file-1", "file-2"})
	jobs, _ := s.Invoke()
	assert.Len(t, jobs, 2)
}
