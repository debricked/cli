package io

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var zipStruct = Zip{}

const fileNameZip = "debricked-test.zip"

func TestNewWriter(t *testing.T) {
	testFile, err := filesystem.Create(fileNameZip)
	defer deleteFile(t, testFile)
	assert.NoError(t, err)

	writer := zipStruct.NewWriter(testFile)
	err = zipStruct.CloseWriter(writer)
	assert.NoError(t, err)
}

func TestFileInfoHeader(t *testing.T) {
	testFile, err := filesystem.Create(fileNameZip)
	defer deleteFile(t, testFile)
	assert.NoError(t, err)

	writer := zipStruct.NewWriter(testFile)
	defer zipStruct.CloseWriter(writer)

	info, _ := filesystem.StatFile(testFile)
	_, err = zipStruct.FileInfoHeader(info)

	assert.NoError(t, err)
}

func TestCreateHeader(t *testing.T) {
	testFile, err := filesystem.Create(fileNameZip)
	defer deleteFile(t, testFile)
	assert.NoError(t, err)

	writer := zipStruct.NewWriter(testFile)
	defer zipStruct.CloseWriter(writer)

	info, _ := filesystem.StatFile(testFile)
	header, _ := zipStruct.FileInfoHeader(info)
	_, err = zipStruct.CreateHeader(writer, header)

	assert.NoError(t, err)
}

func TestDeflate(t *testing.T) {
	deflate := zipStruct.GetDeflate()
	assert.Equal(t, deflate, uint16(8))
}

func TestOpenZip(t *testing.T) {
	r, err := zipStruct.OpenReader("testdata/text.zip")
	assert.NoError(t, err)
	defer zipStruct.CloseReader(r)

	assert.NotNil(t, r, "reader not opened properly")
	_, err = zipStruct.Open(r.File[0])
	assert.NoError(t, err, "should be able to open zip file")
}

func TestCloseZip(t *testing.T) {
	r, err := zipStruct.OpenReader("testdata/text.zip")
	assert.NoError(t, err)
	defer zipStruct.CloseReader(r)

	assert.NotNil(t, r, "reader not opened properly")
	rc, err := zipStruct.Open(r.File[0])
	assert.NoError(t, err, "should be able to open zip file")

	err = zipStruct.Close(rc)
	assert.NoError(t, err, "could not close ReadCloser with zip.Close() properly")
}

func TestOpenCloseReader(t *testing.T) {
	arc := NewArchive(".")
	testFileName := "testdata/ziptest.zip"
	err := arc.ZipFile("testdata/text.txt", testFileName, "ziptest.zip")
	assert.NoError(t, err)

	testFile, err := filesystem.Open(testFileName)
	defer deleteFile(t, testFile)
	assert.NoError(t, err)

	rc, err := zipStruct.OpenReader(testFileName)
	assert.NotNil(t, rc, "OpenReader did not create a viable ReadCloser")
	assert.NoError(t, err, "OpenReader did not properly open ReadCloser")
	err = zipStruct.CloseReader(rc)
	assert.NoError(t, err, "CloseReader did not properly close ReadCloser")
}
