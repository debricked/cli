package testdata

import (
	"io"
	"os"
)

type FileSystemMock struct {
	file          *os.File
	fileInfo      os.FileInfo
	OpenError     error
	StatError     error
	ReadFileError error
	CreateError   error
	RemoveError   error
	StatFileError error
	WriteError    error
}

type IFile interface {
	io.Closer
	io.Reader
	io.ReaderAt
	io.Seeker
	Stat() (os.FileInfo, error)
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

func (fsm FileSystemMock) CloseFile(file *os.File) {
}

func (fsm FileSystemMock) WriteToWriter(_ io.Writer, bytes []byte) (int, error) {
	return 0, fsm.WriteError
}
