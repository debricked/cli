package testdata

import (
	"os"

	"github.com/debricked/cli/internal/io"
)

// FileSystemMock delegates to the real file system, but allows tests to inject
// errors for the operations the policy validator depends on.
type FileSystemMock struct {
	io.FileSystem
	StatError     error
	NotExist      bool
	ReadFileError error
}

func (mock FileSystemMock) Stat(path string) (os.FileInfo, error) {
	if mock.StatError != nil {
		return nil, mock.StatError
	}

	return mock.FileSystem.Stat(path)
}

func (mock FileSystemMock) IsNotExist(err error) bool {
	if mock.StatError != nil {
		return mock.NotExist
	}

	return mock.FileSystem.IsNotExist(err)
}

func (mock FileSystemMock) ReadFile(path string) ([]byte, error) {
	if mock.ReadFileError != nil {
		return nil, mock.ReadFileError
	}

	return mock.FileSystem.ReadFile(path)
}
