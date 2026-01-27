package uv

import (
	"testing"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/stretchr/testify/assert"
)

type jobMock struct {
	job.IJob
}

func TestNewStrategy(t *testing.T) {
	files := []string{"pyproject.toml"}
	strategy := NewStrategy(files)

	assert.Equal(t, files, strategy.files)
}

func TestInvoke(t *testing.T) {
	files := []string{"pyproject.toml"}
	strategy := NewStrategy(files)

	jobs, err := strategy.Invoke()

	assert.NoError(t, err)
	assert.Len(t, jobs, 1)
}
