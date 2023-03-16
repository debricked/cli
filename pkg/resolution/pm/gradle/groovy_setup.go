package gradle

import (
	"embed"

	"github.com/debricked/cli/pkg/resolution/pm/writer"
)

//go:embed gradle-init/gradle-init-script.groovy
var gradleInitScript embed.FS

type ISetupFile interface {
	ReadInitFile() ([]byte, error)
	WriteInitFile() ([]byte, error)
}

type SetupFile struct{}

func (_ SetupFile) ReadInitFile() ([]byte, error) {
	return gradleInitScript.ReadFile("gradle-init/gradle-init-script.groovy")
}

func (sf SetupFile) WriteInitFile(targetFileName string, fileWriter writer.FileWriter) error {
	content, err := sf.ReadInitFile()
	if err != nil {
		return err
	}

	lockFile, err := fileWriter.Create(targetFileName)
	if err != nil {
		return err
	}
	defer lockFile.Close()

	err = fileWriter.Write(lockFile, content)
	return err
}
