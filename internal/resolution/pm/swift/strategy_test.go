package swift

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStrategy(t *testing.T) {
	s := NewStrategy(nil)
	assert.NotNil(t, s)
	assert.Len(t, s.files, 0)

	s = NewStrategy([]string{"file"})
	assert.Len(t, s.files, 1)
}

func TestStrategyInvoke(t *testing.T) {
	cases := [][]string{
		{},
		{"Package.swift"},
		{"a/Package.swift", "b/Package.swift"},
	}

	for _, files := range cases {
		filesCopy := append([]string{}, files...)
		name := "len=" + string(rune(len(filesCopy)))
		t.Run(name, func(t *testing.T) {
			s := NewStrategy(filesCopy)
			jobs, err := s.Invoke()
			assert.NoError(t, err)
			assert.Len(t, jobs, len(filesCopy))
		})
	}
}
