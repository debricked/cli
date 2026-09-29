package testdata

import (
	"embed"
	"io"
	"io/fs"
	"os"

	debio "github.com/debricked/cli/internal/io"
)

// FileSystem is an in-memory debio.IFileSystem test double, keyed by path, giving
// mcpclient tests realistic read/write/exists behavior for config files.
type FileSystem struct {
	Files map[string][]byte
}

var _ debio.IFileSystem = (*FileSystem)(nil)

func NewFileSystem() *FileSystem {
	return &FileSystem{Files: map[string][]byte{}}
}

func (m *FileSystem) ReadFile(path string) ([]byte, error) {
	content, ok := m.Files[path]
	if !ok {
		return nil, os.ErrNotExist
	}

	return content, nil
}

func (m *FileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (m *FileSystem) FsWriteFile(path string, bytes []byte, _ fs.FileMode) error {
	m.Files[path] = append([]byte(nil), bytes...)

	return nil
}

func (m *FileSystem) MkdirAll(_ string, _ fs.FileMode) error {
	return nil
}

func (m *FileSystem) Mkdir(_ string, _ fs.FileMode) error {
	return nil
}

func (m *FileSystem) Open(_ string) (*os.File, error)                   { return nil, nil }
func (m *FileSystem) Create(_ string) (*os.File, error)                 { return nil, nil }
func (m *FileSystem) Stat(_ string) (os.FileInfo, error)                { return nil, nil }
func (m *FileSystem) Remove(_ string) error                             { return nil }
func (m *FileSystem) StatFile(_ *os.File) (os.FileInfo, error)          { return nil, nil }
func (m *FileSystem) CloseFile(_ *os.File)                              {}
func (m *FileSystem) WriteToWriter(_ io.Writer, _ []byte) (int, error)  { return 0, nil }
func (m *FileSystem) MkdirTemp(_ string) (string, error)                { return "", nil }
func (m *FileSystem) RemoveAll(_ string)                                {}
func (m *FileSystem) FsOpenEmbed(_ embed.FS, _ string) (fs.File, error) { return nil, nil }
func (m *FileSystem) FsCloseFile(_ fs.File)                             {}
func (m *FileSystem) FsReadAll(_ fs.File) ([]byte, error)               { return nil, nil }
func (m *FileSystem) Copy(_ io.Writer, _ io.Reader) (int64, error)      { return 0, nil }
