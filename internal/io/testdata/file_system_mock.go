package testdata

import (
	"embed"
	"io"
	"io/fs"
	"os"
)

type FileSystemMock struct {
	file             *os.File
	fileInfo         os.FileInfo
	fsFile           fs.File
	OpenError        error
	StatError        error
	IsNotExistBool   bool
	ReadFileError    error
	CreateError      error
	RemoveError      error
	StatFileError    error
	WriteError       error
	MkdirTempError   error
	FsOpenEmbedError error
	FsReadAllError   error
	FsWriteFileError error
	MkdirError       error
	CopyError        error
	CopySize         int64
}

func (fsm FileSystemMock) Open(path string) (*os.File, error) {
	return fsm.file, fsm.OpenError
}

func (fsm FileSystemMock) Create(path string) (*os.File, error) {
	return fsm.file, fsm.CreateError
}

func (fsm FileSystemMock) Stat(path string) (os.FileInfo, error) {
	return fsm.fileInfo, fsm.StatError
}

func (fsm FileSystemMock) ReadFile(path string) ([]byte, error) {
	return []byte{}, fsm.ReadFileError
}

func (fsm FileSystemMock) Remove(path string) error {
	return fsm.RemoveError
}

func (fsm FileSystemMock) StatFile(file *os.File) (os.FileInfo, error) {
	return fsm.fileInfo, fsm.StatFileError
}

func (fsm FileSystemMock) IsNotExist(err error) bool {
	return fsm.IsNotExistBool
}

func (fsm FileSystemMock) CloseFile(file *os.File) {
}

func (fsm FileSystemMock) WriteToWriter(_ io.Writer, bytes []byte) (int, error) {
	return 0, fsm.WriteError
}

func (fsm FileSystemMock) MkdirTemp(pattern string) (string, error) {
	return pattern, fsm.MkdirTempError
}

func (fsm FileSystemMock) Mkdir(name string, perm fs.FileMode) error {
	return fsm.MkdirError
}

func (fsm FileSystemMock) RemoveAll(path string) {
}

func (fsm FileSystemMock) FsOpenEmbed(file embed.FS, path string) (fs.File, error) {
	return fsm.fsFile, fsm.FsOpenEmbedError
}

func (fsm FileSystemMock) FsCloseFile(file fs.File) {
}

func (fsm FileSystemMock) FsReadAll(file fs.File) ([]byte, error) {
	return []byte{}, fsm.FsReadAllError
}

func (fsm FileSystemMock) FsWriteFile(path string, bytes []byte, perm fs.FileMode) error {
	return fsm.FsWriteFileError
}

func (fsm FileSystemMock) Copy(destination io.Writer, source io.Reader) (int64, error) {
	return fsm.CopySize, fsm.CopyError
}
