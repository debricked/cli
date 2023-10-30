package testdata

import (
	"os"
)

type FileWriterMock struct {
	file      *os.File
	Contents  []byte
	CreateErr error
	WriteErr  error
	CloseErr  error
}

func (fw *FileWriterMock) Create(_ string) (*os.File, error) {
	return fw.file, fw.CreateErr
}

func (fw *FileWriterMock) Write(_ *os.File, bytes []byte) error {
	fw.Contents = append(fw.Contents, bytes...)

	return fw.WriteErr
}

func (fw *FileWriterMock) Close(_ *os.File) error {
	return fw.CloseErr
}
