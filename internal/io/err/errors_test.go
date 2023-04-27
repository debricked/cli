package err

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewErrors(t *testing.T) {
	title := "title"
	errors := NewErrors(title)
	assert.Equal(t, title, errors.title)
	assert.NotNil(t, errors)
	assert.Empty(t, errors.criticalErrs)
	assert.Empty(t, errors.warningErrs)
}

func TestWarning(t *testing.T) {
	errors := NewErrors("")
	warning := fmt.Errorf("error")
	errors.Warning(warning)
	assert.Empty(t, errors.criticalErrs)
	assert.Len(t, errors.warningErrs, 1)
	assert.Contains(t, errors.warningErrs, warning)
}

func TestCritical(t *testing.T) {
	errors := NewErrors("")
	critical := fmt.Errorf("error")
	errors.Critical(critical)
	assert.Empty(t, errors.warningErrs)
	assert.Len(t, errors.criticalErrs, 1)
	assert.Contains(t, errors.criticalErrs, critical)
}

func TestGetWarningErrors(t *testing.T) {
	errors := NewErrors("")
	warning := fmt.Errorf("error")
	errors.Warning(warning)
	assert.Empty(t, errors.GetCriticalErrors())
	assert.Len(t, errors.GetWarningErrors(), 1)
	assert.Contains(t, errors.GetWarningErrors(), warning)
}

func TestGetCriticalErrors(t *testing.T) {
	errors := NewErrors("")
	critical := fmt.Errorf("error")
	errors.Critical(critical)
	assert.Empty(t, errors.GetWarningErrors())
	assert.Len(t, errors.GetCriticalErrors(), 1)
	assert.Contains(t, errors.GetCriticalErrors(), critical)
}

func TestGetAll(t *testing.T) {
	errors := NewErrors("")
	warning := fmt.Errorf("warning")
	critical := fmt.Errorf("critical")
	errors.Warning(warning)
	errors.Critical(critical)
	assert.Len(t, errors.GetAll(), 2)
	assert.Contains(t, errors.GetAll(), warning)
	assert.Contains(t, errors.GetAll(), critical)
}

func TestHasError(t *testing.T) {
	errors := NewErrors("")
	assert.False(t, errors.HasError())

	warning := fmt.Errorf("warning")
	errors.Warning(warning)
	assert.True(t, errors.HasError())

	critical := fmt.Errorf("critical")
	errors.Warning(critical)
	assert.True(t, errors.HasError())
}
