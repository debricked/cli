package java

import (
	"fmt"
	"testing"

	ioFs "github.com/debricked/cli/internal/io"
	ioTestData "github.com/debricked/cli/internal/io/testdata"

	"github.com/stretchr/testify/assert"
)

func TestInitializeSootWrapper(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{}
	tempDir, err := fsMock.MkdirTemp(".tmp")
	assert.NoError(t, err)
	path, err := initializeSootWrapper(fsMock, tempDir)
	assert.NotNil(t, path)
	assert.NoError(t, err)
}

func TestInitializeSootWrapperOpenEmbedError(t *testing.T) {
	errString := "fs open embed error"
	fsMock := ioTestData.FileSystemMock{FsOpenEmbedError: fmt.Errorf(errString)}
	tempDir, err := fsMock.MkdirTemp(".tmp")
	assert.NoError(t, err)
	_, err = initializeSootWrapper(fsMock, tempDir)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errString)
}

func TestInitializeSootWrapperFsReadAllError(t *testing.T) {
	errString := "fs read all error"
	fsMock := ioTestData.FileSystemMock{FsReadAllError: fmt.Errorf(errString)}
	tempDir, err := fsMock.MkdirTemp(".tmp")
	assert.NoError(t, err)
	_, err = initializeSootWrapper(fsMock, tempDir)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errString)
}

func TestInitializeSootWrapperFsWriteFileError(t *testing.T) {
	errString := "fs write file error"
	fsMock := ioTestData.FileSystemMock{FsWriteFileError: fmt.Errorf(errString)}
	tempDir, err := fsMock.MkdirTemp(".tmp")
	assert.NoError(t, err)
	_, err = initializeSootWrapper(fsMock, tempDir)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errString)
}

func TestDownloadSootWrapper(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{}
	arcMock := ioTestData.ArchiveMock{}
	err := downloadSootWrapper(arcMock, fsMock, "soot-wrapper.jar", "11")
	assert.NoError(t, err, "expected no error for downloading soot-wrapper jar")
}

func TestDownloadSootWrapperMkdirTempError(t *testing.T) {
	errString := "mkdir temp error"
	fsMock := ioTestData.FileSystemMock{MkdirTempError: fmt.Errorf(errString)}
	arcMock := ioTestData.ArchiveMock{}
	err := downloadSootWrapper(arcMock, fsMock, "soot-wrapper.jar", "11")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errString)
}

func TestDownloadSootWrapperCreateError(t *testing.T) {
	errString := "create error"
	fsMock := ioTestData.FileSystemMock{CreateError: fmt.Errorf(errString)}
	arcMock := ioTestData.ArchiveMock{}
	err := downloadSootWrapper(arcMock, fsMock, "soot-wrapper.jar", "11")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errString)
}

func TestDownloadSootWrapperUnzipError(t *testing.T) {
	errString := "create error"
	fsMock := ioTestData.FileSystemMock{}
	arcMock := ioTestData.ArchiveMock{UnzipFileError: fmt.Errorf(errString)}
	err := downloadSootWrapper(arcMock, fsMock, "soot-wrapper.jar", "11")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errString)
}

func TestDownloadCompressedSootWrapper(t *testing.T) {
	fs := ioFs.FileSystem{}
	dir, err := fs.MkdirTemp(".test_tmp")
	assert.NoError(t, err, "trying to make temp dir")
	path := dir + "/soot_wrapper.zip"
	file, err := fs.Create(path)
	assert.NoError(t, err, "trying to create file")
	defer file.Close()

	err = downloadCompressedSootWrapper(fs, file, "11")
	assert.NoError(t, err, "expected no error for downloading soot-wrapper")
}

func TestGetSootWrapper(t *testing.T) {
	fs := ioTestData.FileSystemMock{}
	arc := ioTestData.ArchiveMock{}
	sootHandler := SootHandler{}
	tests := []struct {
		name        string
		version     string
		expectError bool
	}{
		{
			name:        "Unsupported version",
			version:     "8",
			expectError: true,
		},
		{
			name:        "Version 11",
			version:     "11",
			expectError: false,
		},
		{
			name:        "Version 17",
			version:     "17",
			expectError: false,
		},
		{
			name:        "Version 21",
			version:     "21",
			expectError: false,
		},
		{
			name:        "Version not int",
			version:     "akjwdm",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := sootHandler.GetSootWrapper(tt.version, fs, arc)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
			}
		})
	}
}

func TestGetSootWrapperDownload(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{StatError: fmt.Errorf(""), IsNotExistBool: true}
	arcMock := ioTestData.ArchiveMock{}
	sootHandler := SootHandler{}
	_, err := sootHandler.GetSootWrapper("17", fsMock, arcMock)
	assert.NoError(t, err)
}

func TestGetSootWrapperInitialize(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{StatError: fmt.Errorf(""), IsNotExistBool: true}
	arcMock := ioTestData.ArchiveMock{}
	sootHandler := SootHandler{}
	_, err := sootHandler.GetSootWrapper("23", fsMock, arcMock)
	assert.NoError(t, err)
}

func TestGetSootWrapperMkdirError(t *testing.T) {
	errString := "mkdir error"
	fsMock := ioTestData.FileSystemMock{MkdirError: fmt.Errorf(errString), StatError: fmt.Errorf(""), IsNotExistBool: true}
	arcMock := ioTestData.ArchiveMock{}
	sootHandler := SootHandler{}
	_, err := sootHandler.GetSootWrapper("11", fsMock, arcMock)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), errString)
}
