package poetry

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
		{"pyproject.toml"},
		{"a/pyproject.toml", "b/pyproject.toml"},
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
