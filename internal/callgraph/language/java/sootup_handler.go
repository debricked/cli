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

// ISootUpHandler is the interface for obtaining the SootUpWrapper JAR path.
// It is the SootUp equivalent of ISootHandler and will be used by the
// dual-engine callgraph runner.
type ISootUpHandler interface {
	GetSootUpWrapper(version string, fs ioFs.IFileSystem, arc ioFs.IArchive) (string, error)
}

// SootUpHandler resolves and extracts the SootUpWrapper.jar.
// cliVersion is used only when a per-version ZIP must be downloaded from
// GitHub Releases; for the embedded JAR (Java 11 bytecode, JVM 11/17/21
// compatible) cliVersion is not needed.
type SootUpHandler struct{ cliVersion string }

// GetSootWrapper keeps compatibility with the existing ISootHandler contract.
// It delegates to GetSootUpWrapper so callers can swap engines without
// changing callgraph/job wiring.
func (sh SootUpHandler) GetSootWrapper(version string, fs ioFs.IFileSystem, arc ioFs.IArchive) (string, error) {
	return sh.GetSootUpWrapper(version, fs, arc)
}

//go:embed embedded/SootUpWrapper.jar
var sootUpJarFS embed.FS

// initializeSootUpWrapper copies the embedded SootUpWrapper.jar into tempDir
// and returns its absolute path.
func (sh SootUpHandler) initializeSootUpWrapper(fs ioFs.IFileSystem, tempDir string) (string, error) {
	jarFile, err := fs.FsOpenEmbed(sootUpJarFS, "embedded/SootUpWrapper.jar")
	if err != nil {
		return "", err
	}
	defer fs.FsCloseFile(jarFile)

	tempJarFile := filepath.Join(tempDir, "SootUpWrapper.jar")

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

// downloadSootUpWrapper downloads a release-specific ZIP from GitHub and
// unzips it to path. Used for future per-version SootUp releases.
func (sh SootUpHandler) downloadSootUpWrapper(arc ioFs.IArchive, fs ioFs.IFileSystem, jarPath string, version string) error {
	dir, err := fs.MkdirTemp(".tmp")
	if err != nil {
		return err
	}

	zipPath := dir + "/sootup_wrapper.zip"
	zipFile, err := fs.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	err = sh.downloadCompressedSootUpWrapper(fs, zipFile, version)
	if err != nil {
		return err
	}

	return arc.UnzipFile(zipPath, jarPath)
}

func (sh SootUpHandler) downloadCompressedSootUpWrapper(fs ioFs.IFileSystem, zipFile *os.File, version string) error {
	fullURLFile := strings.Join([]string{
		"https://github.com/debricked/cli/releases/download/",
		sh.cliVersion,
		"/sootup-wrapper-",
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

// GetSootUpWrapper returns the path to a ready-to-run SootUpWrapper.jar for
// the given Java version.
//
// The SootUp JAR is compiled at Java 11 bytecode level and runs on JVM 11,
// 17, and 21 without modification, so a single embedded artifact covers all
// supported versions.  The JAR is extracted to .debricked/ on first use and
// reused on subsequent runs.
func (sh SootUpHandler) GetSootUpWrapper(version string, fs ioFs.IFileSystem, arc ioFs.IArchive) (string, error) {
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		return "", fmt.Errorf("could not convert version to int")
	}
	if _, err = sh.getSootUpHandlerJavaVersion(versionInt); err != nil {
		return "", err
	}

	debrickedDir := ".debricked"
	if _, err := fs.Stat(debrickedDir); fs.IsNotExist(err) {
		if mkdirErr := fs.Mkdir(debrickedDir, 0755); mkdirErr != nil {
			return "", mkdirErr
		}
	}

	absDebrickedDir, err := filepath.Abs(debrickedDir)
	if err != nil {
		return "", err
	}

	jarPath, err := filepath.Abs(path.Join(debrickedDir, "SootUpWrapper.jar"))
	if err != nil {
		return "", err
	}

	// Always use the embedded JAR; it covers all supported JVM versions.
	if _, err := fs.Stat(jarPath); fs.IsNotExist(err) {
		return sh.initializeSootUpWrapper(fs, absDebrickedDir)
	}

	return jarPath, nil
}

// getSootUpHandlerJavaVersion validates the JVM version and returns the
// canonical version string used by SootUp ("11", "17", or "21").
func (sh SootUpHandler) getSootUpHandlerJavaVersion(version int) (string, error) {
	if version >= 21 {
		return "21", nil
	} else if version >= 17 {
		return "17", nil
	} else if version >= 11 {
		return "11", nil
	}

	return "", fmt.Errorf("lowest supported version for running callgraph generation is 11")
}
