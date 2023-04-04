package gradle

import (
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

type IInitScriptHandler interface {
	ReadInitFile() ([]byte, error)
	WriteInitFile(targetFileName string, fileWriter writer.IFileWriter) error
}

type InitScriptHandler struct{}

func (_ InitScriptHandler) ReadInitFile() ([]byte, error) {
	return gradleInitScript.ReadFile("gradle-init/gradle-init-script.groovy")
}

func (i InitScriptHandler) WriteInitFile(targetFileName string, fileWriter writer.IFileWriter) error {
	content, err := i.ReadInitFile()
	if err != nil {

		return SetupScriptError{message: err.Error()}
	}
	lockFile, err := fileWriter.Create(targetFileName)
	if err != nil {

		return SetupScriptError{message: err.Error()}
	}
	defer lockFile.Close()
	err = fileWriter.Write(lockFile, content)
	if err != nil {

		return SetupScriptError{message: err.Error()}
	}

	return nil
}
