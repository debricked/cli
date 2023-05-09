package finder

import (
	"github.com/debricked/cli/pkg/io/writer"
)

type IInitScriptHandler interface {
	ReadInitFile() ([]byte, error)
	WriteInitFile() error
}

type InitScriptHandler struct {
	groovyScriptPath string
	fileWriter       writer.IFileWriter
}

func (i InitScriptHandler) ReadInitFile() ([]byte, error) {
	return gradleInitScript.ReadFile("embeded/gradle-init-script.groovy")
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
