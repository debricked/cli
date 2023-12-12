package job

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBaseJobError(t *testing.T) {
	error_message := "error"
	jobError := NewBaseJobError(error_message)
	assert.Equal(t, error_message, jobError.err)
	assert.Equal(t, string(""), jobError.documentation)
	assert.NotNil(t, jobError)
}

func TestBaseJobErrorError(t *testing.T) {
	jobError := BaseJobError{
		err:           "error",
		command:       "",
		documentation: "",
		status:        "",
	}
	assert.Equal(t, "error", jobError.Error())
}

func TestBaseJobErrorCommand(t *testing.T) {
	jobError := BaseJobError{
		err:           "",
		command:       "command",
		documentation: "",
		status:        "",
	}
	assert.Equal(t, "command", jobError.Command())
}

func TestBaseJobErrorSetCommand(t *testing.T) {
	jobError := NewBaseJobError("")
	assert.Equal(t, "", jobError.Command())
	jobError.SetCommand("command")
	assert.Equal(t, "command", jobError.Command())
}

func TestBaseJobErrorDocumentation(t *testing.T) {
	jobError := BaseJobError{
		err:           "",
		command:       "",
		documentation: "documentation",
		status:        "",
	}
	assert.Equal(t, "documentation\n", jobError.Documentation())
}

func TestBaseJobErrorSetDocumentation(t *testing.T) {
	jobError := NewBaseJobError("")
	assert.Equal(t, "\n", jobError.Documentation())
	jobError.SetDocumentation("documentation")
	assert.Equal(t, "documentation\n", jobError.Documentation())
}
func TestBaseJobErrorStatus(t *testing.T) {
	jobError := BaseJobError{
		err:           "",
		command:       "",
		documentation: "",
		status:        "status",
	}
	assert.Equal(t, "status", jobError.Status())
}

func TestBaseJobErrorSetStatus(t *testing.T) {
	jobError := NewBaseJobError("")
	assert.Equal(t, "", jobError.Status())
	jobError.SetStatus("status")
	assert.Equal(t, "status", jobError.Status())
}
