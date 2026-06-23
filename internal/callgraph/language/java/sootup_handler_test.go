package java

import (
	"fmt"
	"testing"

	ioTestData "github.com/debricked/cli/internal/io/testdata"

	"github.com/stretchr/testify/assert"
)

var sootUpHandler = SootUpHandler{}

// ── initializeSootUpWrapper ───────────────────────────────────────────────────

func TestInitializeSootUpWrapper(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{}
	tempDir, err := fsMock.MkdirTemp(".tmp")
	assert.NoError(t, err)
	jarPath, err := sootUpHandler.initializeSootUpWrapper(fsMock, tempDir)
	assert.NotEmpty(t, jarPath)
	assert.NoError(t, err)
}

func TestInitializeSootUpWrapperOpenEmbedError(t *testing.T) {
	errString := "fs open embed error"
	fsMock := ioTestData.FileSystemMock{FsOpenEmbedError: fmt.Errorf("%s", errString)} //nolint
	tempDir, err := fsMock.MkdirTemp(".tmp")
	assert.NoError(t, err)
	_, err = sootUpHandler.initializeSootUpWrapper(fsMock, tempDir)
	assert.Error(t, err)
	assert.Equal(t, errString, err.Error())
}

func TestInitializeSootUpWrapperFsReadAllError(t *testing.T) {
	errString := "fs read all error"
	fsMock := ioTestData.FileSystemMock{FsReadAllError: fmt.Errorf("%s", errString)} //nolint
	tempDir, err := fsMock.MkdirTemp(".tmp")
	assert.NoError(t, err)
	_, err = sootUpHandler.initializeSootUpWrapper(fsMock, tempDir)
	assert.Error(t, err)
	assert.Equal(t, errString, err.Error())
}

func TestInitializeSootUpWrapperFsWriteFileError(t *testing.T) {
	errString := "fs write file error"
	fsMock := ioTestData.FileSystemMock{FsWriteFileError: fmt.Errorf("%s", errString)} //nolint
	tempDir, err := fsMock.MkdirTemp(".tmp")
	assert.NoError(t, err)
	_, err = sootUpHandler.initializeSootUpWrapper(fsMock, tempDir)
	assert.Error(t, err)
	assert.Equal(t, errString, err.Error())
}

// ── downloadSootUpWrapper ─────────────────────────────────────────────────────

func TestDownloadSootUpWrapperMkdirTempError(t *testing.T) {
	errString := "mkdir temp error"
	fsMock := ioTestData.FileSystemMock{MkdirTempError: fmt.Errorf("%s", errString)} //nolint
	arcMock := ioTestData.ArchiveMock{}
	err := sootUpHandler.downloadSootUpWrapper(arcMock, fsMock, "sootup-wrapper.jar", "11")
	assert.Error(t, err)
	assert.Equal(t, errString, err.Error())
}

func TestDownloadSootUpWrapperCreateError(t *testing.T) {
	errString := "create error"
	fsMock := ioTestData.FileSystemMock{CreateError: fmt.Errorf("%s", errString)} //nolint
	arcMock := ioTestData.ArchiveMock{}
	err := sootUpHandler.downloadSootUpWrapper(arcMock, fsMock, "sootup-wrapper.jar", "11")
	assert.Error(t, err)
	assert.Equal(t, errString, err.Error())
}

func TestDownloadSootUpWrapperUnzipError(t *testing.T) {
	errString := "unzip error"
	fsMock := ioTestData.FileSystemMock{}
	arcMock := ioTestData.ArchiveMock{UnzipFileError: fmt.Errorf("%s", errString)} //nolint
	err := sootUpHandler.downloadSootUpWrapper(arcMock, fsMock, "sootup-wrapper.jar", "11")
	assert.Error(t, err)
	assert.Equal(t, errString, err.Error())
}

// ── getSootUpHandlerJavaVersion ───────────────────────────────────────────────

func TestGetSootUpHandlerJavaVersionTooLow(t *testing.T) {
	_, err := sootUpHandler.getSootUpHandlerJavaVersion(8)
	assert.Error(t, err)
}

func TestGetSootUpHandlerJavaVersion11(t *testing.T) {
	v, err := sootUpHandler.getSootUpHandlerJavaVersion(11)
	assert.NoError(t, err)
	assert.Equal(t, "11", v)
}

func TestGetSootUpHandlerJavaVersion17(t *testing.T) {
	v, err := sootUpHandler.getSootUpHandlerJavaVersion(17)
	assert.NoError(t, err)
	assert.Equal(t, "17", v)
}

func TestGetSootUpHandlerJavaVersion21(t *testing.T) {
	v, err := sootUpHandler.getSootUpHandlerJavaVersion(21)
	assert.NoError(t, err)
	assert.Equal(t, "21", v)
}

func TestGetSootUpHandlerJavaVersionAbove21(t *testing.T) {
	v, err := sootUpHandler.getSootUpHandlerJavaVersion(23)
	assert.NoError(t, err)
	assert.Equal(t, "21", v)
}

// Component tests moved to sootup_handler_component_test.go for consistency with soot_handler structure
