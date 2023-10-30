package io

import (
	"os"
)

type IFileWriter interface {
	Write(file *os.File, p []byte) error
	Create(name string) (*os.File, error)
	Close(file *os.File) error
}

type FileWriter struct{}

func (fw FileWriter) Create(name string) (*os.File, error) {
	return os.Create(name)
}

func (fw FileWriter) Write(file *os.File, p []byte) error {
	_, err := file.Write(p)

	return err
}

func (fw FileWriter) Close(file *os.File) error {
	return file.Close()
}
