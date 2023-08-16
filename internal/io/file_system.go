package io

import (
	"embed"
	"io"
	"io/fs"
	"os"
)

type IFileSystem interface {
	Open(path string) (*os.File, error)
	Create(path string) (*os.File, error)
	Stat(path string) (os.FileInfo, error)
	ReadFile(path string) ([]byte, error)
	Remove(path string) error
	StatFile(file *os.File) (os.FileInfo, error)
	IsNotExist(err error) bool
	CloseFile(file *os.File)
	WriteToWriter(writer io.Writer, content []byte) (int, error)
	MkdirTemp(pattern string) (string, error)
	RemoveAll(path string)
	FsOpenEmbed(file embed.FS, path string) (fs.File, error)
	FsCloseFile(file fs.File)
	FsReadAll(file fs.File) ([]byte, error)
	FsWriteFile(path string, bytes []byte, perm fs.FileMode) error
}

type FileSystem struct{}

func (_ FileSystem) Open(path string) (*os.File, error) {
	return os.Open(path)
}

func (_ FileSystem) Create(path string) (*os.File, error) {
	return os.Create(path)
}

func (_ FileSystem) Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (_ FileSystem) IsNotExist(err error) bool {
	return os.IsNotExist(err)
}

func (_ FileSystem) StatFile(file *os.File) (os.FileInfo, error) {
	return file.Stat()
}

func (_ FileSystem) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (_ FileSystem) Remove(path string) error {
	return os.Remove(path)
}

func (_ FileSystem) CloseFile(file *os.File) {
	file.Close()
}

func (_ FileSystem) WriteToWriter(writer io.Writer, content []byte) (int, error) {
	return writer.Write(content)
}

func (_ FileSystem) MkdirTemp(pattern string) (string, error) {
	return os.MkdirTemp("", pattern)
}

func (_ FileSystem) RemoveAll(path string) {
	os.RemoveAll(path)
}

func (_ FileSystem) FsOpenEmbed(file embed.FS, path string) (fs.File, error) {
	return file.Open(path)
}

func (_ FileSystem) FsCloseFile(file fs.File) {
	file.Close()
}

func (_ FileSystem) FsReadAll(file fs.File) ([]byte, error) {
	return io.ReadAll(file)
}

func (_ FileSystem) FsWriteFile(path string, bytes []byte, perm fs.FileMode) error {
	return os.WriteFile(path, bytes, perm)
}
