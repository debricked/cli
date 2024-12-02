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
	ReaderCloser      *zip.ReadCloser
	FileHeaderError   error
	CreateHeaderError error
	CloseWriterError  error
	OpenReaderError   error
	CloseReaderError  error
	OpenError         error
	CloseError        error
	ReadError         error
	ReadCloser        io.ReadCloser
}

func (zm ZipMock) NewWriter(file *os.File) *zip.Writer {
	return zm.writer
}

func (zm ZipMock) FileInfoHeader(fileInfo fs.FileInfo) (*zip.FileHeader, error) {
	return &zip.FileHeader{}, zm.FileHeaderError
}

func (zm ZipMock) OpenReader(source string) (*zip.ReadCloser, error) {

	return zm.ReaderCloser, zm.OpenReaderError
}

func (zm ZipMock) CloseReader(reader *zip.ReadCloser) error {
	return zm.CloseReaderError
}

func (zm ZipMock) GetDeflate() uint16 {
	return zip.Deflate
}

func (zm ZipMock) CreateHeader(writer *zip.Writer, header *zip.FileHeader) (io.Writer, error) {
	return zm.createHeader, zm.CreateHeaderError
}

func (zm ZipMock) CloseWriter(writer *zip.Writer) error {
	return zm.CloseWriterError
}

func (zm ZipMock) Open(file *zip.File) (io.ReadCloser, error) {
	return zm.ReadCloser, zm.OpenError
}

func (zm ZipMock) Close(rc io.ReadCloser) error {
	return zm.CloseError
}

type MockReader struct {
	ReadBytes int
	ReadError error
}

func (r *MockReader) Read(p []byte) (int, error) {
	return r.ReadBytes, r.ReadError
}

func MockFileSlice(contents []string) []*zip.File {
	files := []*zip.File{}
	return files
}
