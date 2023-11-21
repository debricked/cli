package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewError(t *testing.T) {
	error_message := "error"
	jobError := NewPMJobError(error_message)
	assert.Equal(t, error_message, jobError.err)
	assert.Equal(t, string(""), jobError.cmd)
	assert.Equal(
		t,
		string("No specific documentation for this problem yet, please report it to us! :)"),
		jobError.doc,
	)
	assert.Equal(t, string(""), jobError.status)
	assert.NotNil(t, jobError)
}

func TestPMJobErrorError(t *testing.T) {
	jobError := PMJobError{
		err:    "error",
		cmd:    "",
		doc:    "",
		status: "",
	}
	assert.Equal(t, "error", jobError.Error())
}

func TestPMJobErrorCommand(t *testing.T) {
	jobError := PMJobError{
		err:    "",
		cmd:    "command",
		doc:    "",
		status: "",
	}
	assert.Equal(t, "`command`\n", jobError.Command())
}

func TestPMJobErrorSetCommand(t *testing.T) {
	jobError := NewPMJobError("")
	assert.Equal(t, "", jobError.Command())
	jobError.SetCommand("command")
	assert.Equal(t, "`command`\n", jobError.Command())
}

func TestPMJobErrorDocumentation(t *testing.T) {
	jobError := PMJobError{
		err:    "",
		cmd:    "",
		doc:    "documentation",
		status: "",
	}
	assert.Equal(t, "documentation\n", jobError.Documentation())
}

func TestPMJobErrorSetDocumentation(t *testing.T) {
	jobError := NewPMJobError("")
	assert.Equal(
		t,
		string("No specific documentation for this problem yet, please report it to us! :)\n"),
		jobError.Documentation(),
	)
	jobError.SetDocumentation("documentation")
	assert.Equal(t, "documentation\n", jobError.Documentation())
}

func TestPMJobErrorStatus(t *testing.T) {
	jobError := PMJobError{
		err:    "",
		cmd:    "",
		doc:    "",
		status: "status",
	}
	assert.Equal(t, "status", jobError.Status())
}

func TestPMJobErrorSetStatus(t *testing.T) {
	jobError := NewPMJobError("")
	assert.Equal(t, "", jobError.Status())
	jobError.SetStatus("status")
	assert.Equal(t, "status", jobError.Status())
}
