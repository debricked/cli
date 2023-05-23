package io

import (
	"io"
	"os"
)

type IFileSystem interface {
	Open(path string) (*os.File, error)
	Create(path string) (*os.File, error)
	Stat(path string) (os.FileInfo, error)
	ReadFile(path string) ([]byte, error)
	Remove(path string) error
	StatFile(file *os.File) (os.FileInfo, error)
	CloseFile(file *os.File)
	WriteToWriter(writer io.Writer, content []byte) (int, error)
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
