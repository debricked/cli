package io

import (
	"fmt"
	"testing"

	ioTestData "github.com/debricked/cli/internal/io/testdata"
	"github.com/stretchr/testify/assert"
)

func TestZipWMock(t *testing.T) {
	fsMock := ioTestData.FileSystemMock{}
	zipMock := ioTestData.ZipMock{}
	a := Archive{
		workingDirectory: "nonexisting",
		fs:               fsMock,
		zip:              zipMock,
	}

	err := a.ZipFile("testdir", "targettest")
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

	err := a.ZipFile("testdir", "targettest")
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

	err := a.ZipFile("testdir", "targettest")
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

	err := a.ZipFile("testdir", "targettest")
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

	err := a.ZipFile("testdir", "targettest")
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

	err := a.ZipFile("testdir", "targettest")
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

	err := a.ZipFile("testdir", "targettest")
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
