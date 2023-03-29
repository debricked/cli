package gradle

import (
	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

type IInitFileHandler interface {
	ReadInitFile() ([]byte, error)
	WriteInitFile(targetFileName string, fileWriter writer.IFileWriter) error
}

type InitFileHandler struct{}

func (_ InitFileHandler) ReadInitFile() ([]byte, error) {
	return gradleInitScript.ReadFile("gradle-init/gradle-init-script.groovy")
}

func (i InitFileHandler) WriteInitFile(targetFileName string, fileWriter writer.IFileWriter) error {
	content, err := i.ReadInitFile()
	if err != nil {
		return GradleSetupScriptError{message: err.Error()}
	}
	lockFile, err := fileWriter.Create(targetFileName)
	if err != nil {
		return GradleSetupScriptError{message: err.Error()}
	}
	defer lockFile.Close()
	err = fileWriter.Write(lockFile, content)
	if err != nil {
		return GradleSetupScriptError{message: err.Error()}
	}
	return nil

}
