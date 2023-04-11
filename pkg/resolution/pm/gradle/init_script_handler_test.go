package gradle

import (
	"embed"
	"errors"
	"testing"

	writerTestdata "github.com/debricked/cli/pkg/resolution/pm/writer/testdata"
	"github.com/stretchr/testify/assert"
)

func TestWriteInitFile(t *testing.T) {
	createErr := errors.New("create-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}

	sf := InitScriptHandler{}
	err := sf.WriteInitFile("file", fileWriterMock)
	assert.Equal(t, SetupScriptError{createErr.Error()}, err)

	fileWriterMock = &writerTestdata.FileWriterMock{WriteErr: createErr}
	err = sf.WriteInitFile("file", fileWriterMock)
	assert.Equal(t, SetupScriptError{createErr.Error()}, err)
}

func TestWriteInitFileNoInitFile(t *testing.T) {
	sf := InitScriptHandler{}
	oldGradleInitScript := gradleInitScript
	defer func() {
		gradleInitScript = oldGradleInitScript
	}()
	gradleInitScript = embed.FS{}
	err := sf.WriteInitFile("file", nil)
	readErr := errors.New("open gradle-init/gradle-init-script.groovy: file does not exist")
	assert.Equal(t, SetupScriptError{readErr.Error()}, err)

}
