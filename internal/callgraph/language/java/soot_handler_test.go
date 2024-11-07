package java

import (
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

func TestDownloadSootWrapper(t *testing.T) {
	fs := ioFs.FileSystem{}
	dir, _ := fs.MkdirTemp(".test_tmp")
	zip := ioFs.Zip{}
	arc := ioFs.NewArchiveWithStructs(dir, fs, zip)
	err := downloadSootWrapper(arc, fs, "soot-wrapper.jar", "11")
	assert.NoError(t, err, "expected no error for downloading soot-wrapper jar")
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
