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

type SootHandler struct{ cliVersion string }

//go:embed embedded/SootWrapper.jar
var jarCallGraph embed.FS

func (sh SootHandler) initializeSootWrapper(fs ioFs.IFileSystem, tempDir string) (string, error) {
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

func (sh SootHandler) downloadSootWrapper(arc ioFs.IArchive, fs ioFs.IFileSystem, path string, version string) error {
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

	err = sh.downloadCompressedSootWrapper(fs, zipFile, version)
	if err != nil {

		return err
	}

	return arc.UnzipFile(zipPath, path)
}

func (sh SootHandler) downloadCompressedSootWrapper(fs ioFs.IFileSystem, zipFile *os.File, version string) error {
	fullURLFile := strings.Join([]string{
		"https://github.com/debricked/cli/releases/download/",
		sh.cliVersion,
		"/soot-wrapper-",
		version,
		".zip",
	}, "")
	fmt.Println("URL=", fullURLFile)

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
		return "", fmt.Errorf("could not convert version to int")
	}
	version, err = sh.getSootHandlerJavaVersion(versionInt)
	if err != nil {
		return "", err
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
		if version == "21" {
			return sh.initializeSootWrapper(fs, debrickedDir)
		}

		return path, sh.downloadSootWrapper(arc, fs, path, version)
	}

	return path, nil
}

func (sh SootHandler) getSootHandlerJavaVersion(version int) (string, error) {
	if version >= 21 {
		return "21", nil
	} else if version >= 17 {
		return "17", nil
	} else if version >= 11 {
		return "11", nil
	} else {
		return "", fmt.Errorf("lowest supported version for running callgraph generation is 11")
	}
}
