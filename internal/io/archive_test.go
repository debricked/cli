package io

import (
	"archive/zip"
	"fmt"
	"testing"

	ioTestData "github.com/debricked/cli/internal/io/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewArchive(t *testing.T) {
	archive := NewArchive(".")
	assert.NotNil(t, archive.fs)
	assert.NotNil(t, archive.zip)

}

func TestNewArchiveWithStructs(t *testing.T) {
	archive := NewArchiveWithStructs(".", FileSystem{}, Zip{})
	assert.NotNil(t, archive.fs)
	assert.NotNil(t, archive.zip)

}

func TestZipWMock(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{}
	zipMock := ioTestData.ZipMock{}
	a := Archive{
		workingDirectory: "nonexisting",
		fs:               fsMock,
		zip:              zipMock,
	}

	err := a.ZipFile("testdir", "targettest", "zippedName")
	assert.Nil(t, err)
}

func TestZipReadFileError(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{ReadFileError: fmt.Errorf("error")}
	zipMock := ioTestData.ZipMock{}
	a := Archive{
		workingDirectory: "nonexisting",
		fs:               fsMock,
		zip:              zipMock,
	}

	err := a.ZipFile("testdir", "targettest", "zippedName")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "error")
}

func TestZipCreateError(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{CreateError: fmt.Errorf("error")}
	zipMock := ioTestData.ZipMock{}
	a := Archive{
		workingDirectory: "nonexisting",
		fs:               fsMock,
		zip:              zipMock,
	}

	err := a.ZipFile("testdir", "targettest", "zippedName")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "error")
}

func TestStatFileError(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{StatFileError: fmt.Errorf("error")}
	zipMock := ioTestData.ZipMock{}
	a := Archive{
		workingDirectory: "nonexisting",
		fs:               fsMock,
		zip:              zipMock,
	}

	err := a.ZipFile("testdir", "targettest", "zippedName")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "error")
}

func TestFileInfoError(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{}
	zipMock := ioTestData.ZipMock{FileHeaderError: fmt.Errorf("error")}
	a := Archive{
		workingDirectory: "nonexisting",
		fs:               fsMock,
		zip:              zipMock,
	}

	err := a.ZipFile("testdir", "targettest", "zippedName")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "error")
}

func TestCreateHeaderError(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{}
	zipMock := ioTestData.ZipMock{CreateHeaderError: fmt.Errorf("error")}
	a := Archive{
		workingDirectory: "nonexisting",
		fs:               fsMock,
		zip:              zipMock,
	}

	err := a.ZipFile("testdir", "targettest", "zippedName")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "error")
}

func TestWriteToWriterError(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{WriteError: fmt.Errorf("error")}
	zipMock := ioTestData.ZipMock{}
	a := Archive{
		workingDirectory: "nonexisting",
		fs:               fsMock,
		zip:              zipMock,
	}

	err := a.ZipFile("testdir", "targettest", "zippedName")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "error")
}

func TestUnzipFileError(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{}
	zipMock := ioTestData.ZipMock{OpenReaderError: fmt.Errorf("error")}
	a := Archive{
		workingDirectory: "nonexisting",
		fs:               fsMock,
		zip:              zipMock,
	}
	err := a.UnzipFile("test_source_path", "test_target_path")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "error")
}

func TestB64(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{}
	a := Archive{
		workingDirectory: "nonexisting",
		fs:               fsMock,
	}

	err := a.B64("testdir", "targettest")
	assert.Nil(t, err)
}

func TestB64ReadFileError(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{ReadFileError: fmt.Errorf("error")}
	a := Archive{
		workingDirectory: "nonexisting",
		fs:               fsMock,
	}

	err := a.B64("testdir", "targettest")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "error")
}

func TestB64CreateError(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{CreateError: fmt.Errorf("error")}
	a := Archive{
		workingDirectory: "nonexisting",
		fs:               fsMock,
	}

	err := a.B64("testdir", "targettest")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "error")
}

func TestB64WriteError(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{WriteError: fmt.Errorf("error")}
	a := Archive{
		workingDirectory: "nonexisting",
		fs:               fsMock,
	}

	err := a.B64("testdir", "targettest")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "error")
}

func TestCleanup(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{WriteError: fmt.Errorf("error")}
	a := Archive{
		workingDirectory: "nonexisting",
		fs:               fsMock,
	}

	err := a.Cleanup("testdir")
	assert.Nil(t, err)
}

func TestUnzipFileReadFileError(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{}
	zipMock := ioTestData.ZipMock{OpenReaderError: fmt.Errorf("%s", t.Name())}
	a := Archive{
		workingDirectory: ".",
		fs:               fsMock,
		zip:              zipMock,
	}
	err := a.UnzipFile("", "")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), t.Name())
}

func TestUnzipFileCreateError(t *testing.T) {
	reader := zip.Reader{
		File: []*zip.File{nil},
	}
	readCloser := zip.ReadCloser{Reader: reader} //nolint
	fsMock := ioTestData.FileSystemMock{CreateError: fmt.Errorf("%s", t.Name())}
	zipMock := ioTestData.ZipMock{ReaderCloser: &readCloser}
	a := Archive{
		workingDirectory: ".",
		fs:               fsMock,
		zip:              zipMock,
	}
	err := a.UnzipFile("", "")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), t.Name())
}

func TestUnzipFileOpenError(t *testing.T) {
	reader := zip.Reader{
		File: []*zip.File{nil},
	}
	readCloser := zip.ReadCloser{Reader: reader} //nolint
	fsMock := ioTestData.FileSystemMock{}
	zipMock := ioTestData.ZipMock{ReaderCloser: &readCloser, OpenError: fmt.Errorf("%s", t.Name())}
	a := Archive{
		workingDirectory: ".",
		fs:               fsMock,
		zip:              zipMock,
	}
	err := a.UnzipFile("", "")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), t.Name())
}
func TestUnzipFileCopyError(t *testing.T) {
	r, err := zipStruct.OpenReader("testdata/text.zip")
	assert.NoError(t, err)
	defer zipStruct.CloseReader(r) //nolint

	fsMock := ioTestData.FileSystemMock{CopyError: fmt.Errorf("%s", t.Name())}
	zipMock := ioTestData.ZipMock{ReaderCloser: r}
	a := Archive{
		workingDirectory: ".",
		fs:               fsMock,
		zip:              zipMock,
	}
	err = a.UnzipFile("", "")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), t.Name())
}

func TestUnzipFileNotSingleFile(t *testing.T) {
	reader := zip.Reader{
		File: []*zip.File{nil, nil},
	}
	readCloser := zip.ReadCloser{Reader: reader} //nolint
	defer zipStruct.CloseReader(&readCloser)     //nolint

	fsMock := ioTestData.FileSystemMock{}
	zipMock := ioTestData.ZipMock{ReaderCloser: &readCloser}
	a := Archive{
		workingDirectory: ".",
		fs:               fsMock,
		zip:              zipMock,
	}
	err := a.UnzipFile("", "")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "cannot unzip archive which does not contain exactly one file")
}
