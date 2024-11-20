package java

import (
	"embed"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	ioFs "github.com/debricked/cli/internal/io"
)

type ISootHandler interface {
	GetSootWrapper(version string, fs ioFs.IFileSystem, arc ioFs.IArchive) (string, error)
}

type SootHandler struct{}

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

func downloadSootWrapper(arc ioFs.IArchive, fs ioFs.IFileSystem, path string, version string) error {
	dir, err := fs.MkdirTemp(".tmp")
	if err != nil {

		return err
	}

	zipPath := dir + "/soot_wrapper.zip"
	zipFile, err := fs.Create(zipPath)
	if err != nil {

		return err
	}
	defer zipFile.Close()

	err = downloadCompressedSootWrapper(fs, zipFile, version)
	if err != nil {

		return err
	}

	return arc.UnzipFile(zipPath, path)
}

func downloadCompressedSootWrapper(fs ioFs.IFileSystem, zipFile *os.File, version string) error {
	fullURLFile := strings.Join([]string{
		"https://github.com/debricked/cli/releases/download/v2.2.0/soot-wrapper-",
		version,
		".zip",
	}, "")

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path

			return nil
		},
	}
	resp, err := client.Get(fullURLFile)
	if err != nil {

		return err
	}
	defer resp.Body.Close()

	_, err = fs.Copy(zipFile, resp.Body)

	return err
}

func (sh SootHandler) GetSootWrapper(version string, fs ioFs.IFileSystem, arc ioFs.IArchive) (string, error) {
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		return "", fmt.Errorf("error when trying to convert java version string to int")
	}
	if versionInt < 11 {
		return "", fmt.Errorf("lowest supported version for running callgraph generation is 11")
	}
	debrickedDir := ".debricked"
	if _, err := fs.Stat(debrickedDir); fs.IsNotExist(err) {
		err := fs.Mkdir(debrickedDir, 0755)
		if err != nil {
			return "", err
		}
	}
	path, err := filepath.Abs(path.Join(debrickedDir, "soot-wrapper.jar"))
	if err != nil {

		return "", err
	}
	if _, err := fs.Stat(path); fs.IsNotExist(err) {
		if versionInt >= 21 {
			return initializeSootWrapper(fs, debrickedDir)
		}
		if versionInt >= 17 {
			version = "17"
		} else {
			version = "11"
		} // Handling correct jar to install

		return path, downloadSootWrapper(arc, fs, path, version)
	}

	return path, nil
}
