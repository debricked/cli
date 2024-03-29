package testdata

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
)

type ZipMock struct {
	writer *zip.Writer
	// fileHeader      *zip.FileHeader
	createHeader      io.Writer
	FileHeaderError   error
	CreateHeaderError error
	CloseError        error
}

func (zm ZipMock) NewWriter(file *os.File) *zip.Writer {
	return zm.writer
}

func (zm ZipMock) FileInfoHeader(fileInfo fs.FileInfo) (*zip.FileHeader, error) {
	return &zip.FileHeader{}, zm.FileHeaderError
}

func (_ ZipMock) GetDeflate() uint16 {
	return zip.Deflate
}

func (zm ZipMock) CreateHeader(writer *zip.Writer, header *zip.FileHeader) (io.Writer, error) {
	return zm.createHeader, zm.CreateHeaderError
}

func (zm ZipMock) Close(writer *zip.Writer) error {
	return zm.CloseError
}
