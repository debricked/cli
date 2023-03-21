package gradle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStrategy(t *testing.T) {
	s := NewStrategy(nil, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 0)

	s = NewStrategy([]string{}, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 0)

	s = NewStrategy([]string{"file"}, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 1)

	s = NewStrategy([]string{"file-1", "file-2"}, nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 2)
}

func TestInvokeNoFiles(t *testing.T) {
	s := NewStrategy([]string{}, nil)
	jobs, _ := s.Invoke()
	assert.Empty(t, jobs)
}

func TestInvokeOneFile(t *testing.T) {
	s := NewStrategy([]string{"file"}, nil)
	jobs, _ := s.Invoke()
	assert.Len(t, jobs, 1)
}

func TestInvokeManyFiles(t *testing.T) {
	s := NewStrategy([]string{"file-1", "file-2"}, nil)
	jobs, _ := s.Invoke()
	assert.Len(t, jobs, 2)
}
