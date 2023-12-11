package job

import (
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
	warning := NewBaseJobError("error")
	errors.Warning(warning)
	assert.Empty(t, errors.criticalErrs)
	assert.Len(t, errors.warningErrs, 1)
	assert.Contains(t, errors.warningErrs, warning)
}

func TestCritical(t *testing.T) {
	errors := NewErrors("")
	critical := NewBaseJobError("error")
	errors.Critical(critical)
	assert.Empty(t, errors.warningErrs)
	assert.Len(t, errors.criticalErrs, 1)
	assert.Contains(t, errors.criticalErrs, critical)
}

func TestAppend(t *testing.T) {
	errors := NewErrors("")

	critical1 := NewBaseJobError("critical")
	critical1.SetIsCritical(true)
	errors.Append(critical1)

	critical2 := NewBaseJobError("another critical")
	errors.Append(critical2)

	warning := NewBaseJobError("warning")
	warning.SetIsCritical(false)
	errors.Append(warning)

	assert.Len(t, errors.warningErrs, 1)
	assert.Len(t, errors.criticalErrs, 2)
	assert.Contains(t, errors.criticalErrs, critical1)
	assert.Contains(t, errors.criticalErrs, critical2)
	assert.Contains(t, errors.warningErrs, warning)
}

func TestGetWarningErrors(t *testing.T) {
	errors := NewErrors("")
	warning := NewBaseJobError("error")
	errors.Warning(warning)
	assert.Empty(t, errors.GetCriticalErrors())
	assert.Len(t, errors.GetWarningErrors(), 1)
	assert.Contains(t, errors.GetWarningErrors(), warning)
}

func TestGetCriticalErrors(t *testing.T) {
	errors := NewErrors("")
	critical := NewBaseJobError("critical")
	errors.Critical(critical)
	assert.Empty(t, errors.GetWarningErrors())
	assert.Len(t, errors.GetCriticalErrors(), 1)
	assert.Contains(t, errors.GetCriticalErrors(), critical)
}

func TestGetAll(t *testing.T) {
	errors := NewErrors("")
	warning := NewBaseJobError("warning")
	critical := NewBaseJobError("critical")
	errors.Warning(warning)
	errors.Critical(critical)
	assert.Len(t, errors.GetAll(), 2)
	assert.Contains(t, errors.GetAll(), warning)
	assert.Contains(t, errors.GetAll(), critical)
}

func TestHasError(t *testing.T) {
	errors := NewErrors("")
	assert.False(t, errors.HasError())

	warning := NewBaseJobError("warning")
	errors.Warning(warning)
	assert.True(t, errors.HasError())
	critical := NewBaseJobError("critical")
	errors.Warning(critical)
	assert.True(t, errors.HasError())
}
