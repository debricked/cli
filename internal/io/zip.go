package io

import (
	"archive/zip"
	"io"
	"io/fs"
	"os"
)

type IZip interface {
	NewWriter(file *os.File) *zip.Writer
	FileInfoHeader(fileInfo fs.FileInfo) (*zip.FileHeader, error)
	GetDeflate() uint16
	CreateHeader(writer *zip.Writer, header *zip.FileHeader) (io.Writer, error)
	Close(writer *zip.Writer) error
	OpenReader(source string) (*zip.ReadCloser, error)
	CloseReader(reader *zip.ReadCloser) error
}

type Zip struct{}

func (_ Zip) NewWriter(file *os.File) *zip.Writer {
	return zip.NewWriter(file)
}

func (_ Zip) FileInfoHeader(fileInfo fs.FileInfo) (*zip.FileHeader, error) {
	return zip.FileInfoHeader(fileInfo)
}

func (_ Zip) GetDeflate() uint16 {
	return zip.Deflate
}

func (z Zip) CreateHeader(writer *zip.Writer, header *zip.FileHeader) (io.Writer, error) {
	return writer.CreateHeader(header)
}

func (z Zip) Close(writer *zip.Writer) error {
	return writer.Close()
}

func (z Zip) OpenReader(source string) (*zip.ReadCloser, error) {
	return zip.OpenReader(source)
}

func (z Zip) CloseReader(reader *zip.ReadCloser) error {
	return reader.Close()
}
