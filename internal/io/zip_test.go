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
	err = zipStruct.Close(writer)
	assert.NoError(t, err)
}

func TestFileInfoHeader(t *testing.T) {
	testFile, err := filesystem.Create(fileNameZip)
	defer deleteFile(t, testFile)
	assert.NoError(t, err)

	writer := zipStruct.NewWriter(testFile)
	defer zipStruct.Close(writer)

	info, _ := filesystem.StatFile(testFile)
	_, err = zipStruct.FileInfoHeader(info)

	assert.NoError(t, err)
}

func TestCreateHeader(t *testing.T) {
	testFile, err := filesystem.Create(fileNameZip)
	defer deleteFile(t, testFile)
	assert.NoError(t, err)

	writer := zipStruct.NewWriter(testFile)
	defer zipStruct.Close(writer)

	info, _ := filesystem.StatFile(testFile)
	header, _ := zipStruct.FileInfoHeader(info)
	_, err = zipStruct.CreateHeader(writer, header)

	assert.NoError(t, err)
}

func TestDeflate(t *testing.T) {
	deflate := zipStruct.GetDeflate()
	assert.Equal(t, deflate, uint16(8))
}
