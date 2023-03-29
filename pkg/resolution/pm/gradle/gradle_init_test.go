package gradle

import (
	"errors"
	"testing"

	writerTestdata "github.com/debricked/cli/pkg/resolution/pm/writer/testdata"
	"github.com/stretchr/testify/assert"
)

func TestWriteInitFile(t *testing.T) {

	// test failing from writer
	createErr := errors.New("create-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}

	sf := InitFileHandler{}
	err := sf.WriteInitFile("file", fileWriterMock)
	assert.Equal(t, GradleSetupScriptError{createErr.Error()}, err)

	fileWriterMock = &writerTestdata.FileWriterMock{WriteErr: createErr}
	err = sf.WriteInitFile("file", fileWriterMock)
	assert.Equal(t, GradleSetupScriptError{createErr.Error()}, err)

}
