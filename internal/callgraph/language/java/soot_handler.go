package java

import (
	"embed"
	"path/filepath"

	ioFs "github.com/debricked/cli/internal/io"
)

//go:embed embedded/SootWrapper.jar
var jarCallGraph embed.FS

func initializeSootWrapper(fs ioFs.IFileSystem, tempDir string) (string, error) {
	jarFile, err := fs.FsOpenEmbed(jarCallGraph, "embedded/SootWrapper.jar")
	if err != nil {
		return "", err
	}
	defer fs.FsCloseFile(jarFile)

	tempJarFile := filepath.Join(tempDir, "SootWrapper.jar")

	jarBytes, err := fs.FsReadAll(jarFile)
	if err != nil {

		return "", err
	}

	err = fs.FsWriteFile(tempJarFile, jarBytes, 0600)
	if err != nil {

		return "", err
	}

	return tempJarFile, nil
}
