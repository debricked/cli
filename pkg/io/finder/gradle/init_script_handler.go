package gradle

import (
	"github.com/debricked/cli/pkg/io/writer"
)

type IInitScriptHandler interface {
	ReadInitFile() ([]byte, error)
	WriteInitFile() error
}

type InitScriptHandler struct {
	groovyScriptPath string
	initPath         string
	fileWriter       writer.IFileWriter
}

func NewScriptHandler(groovyScriptPath string, initPath string, fileWriter writer.IFileWriter) InitScriptHandler {
	return InitScriptHandler{
		groovyScriptPath,
		initPath,
		fileWriter,
	}
}

func (i InitScriptHandler) ReadInitFile() ([]byte, error) {
	return gradleInitScript.ReadFile(i.initPath)
}

func (i InitScriptHandler) WriteInitFile() error {
	content, err := i.ReadInitFile()
	if err != nil {

		return SetupScriptError{message: err.Error()}
	}
	lockFile, err := i.fileWriter.Create(i.groovyScriptPath)
	if err != nil {

		return SetupScriptError{message: err.Error()}
	}
	defer lockFile.Close()
	err = i.fileWriter.Write(lockFile, content)
	if err != nil {

		return SetupScriptError{message: err.Error()}
	}

	return nil
}
