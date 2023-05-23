package io

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var filesystem = FileSystem{}

const fileNameFS = "debricked-test-fs"

func TestCreateFile(t *testing.T) {
	testFile, err := filesystem.Create(fileNameFS)
	assert.NoError(t, err)
	assert.NotNil(t, testFile)
	defer deleteFile(t, testFile)
}

func TestOpen(t *testing.T) {
	testFile, _ := filesystem.Create(fileNameFS)
	defer deleteFile(t, testFile)

	_, err := filesystem.Open(fileNameFS)

	assert.NoError(t, err)
}

func TestStat(t *testing.T) {
	_, err := filesystem.Stat(fileNameFS)
	assert.NotNil(t, err)

	testFile, _ := filesystem.Create(fileNameFS)
	defer deleteFile(t, testFile)

	_, err = filesystem.Stat(fileNameFS)

	assert.NoError(t, err)
}

func TestStatFile(t *testing.T) {
	testFile, _ := filesystem.Create(fileNameFS)
	defer deleteFile(t, testFile)

	_, err := filesystem.StatFile(testFile)

	assert.NoError(t, err)
}

func TestReadFile(t *testing.T) {
	testFile, _ := filesystem.Create(fileNameFS)
	defer deleteFile(t, testFile)

	_, err := filesystem.ReadFile(fileNameFS)

	assert.NoError(t, err)
}

func TestRemoveFile(t *testing.T) {
	testFile, _ := filesystem.Create(fileNameFS)
	defer testFile.Close()

	_, err := filesystem.Stat(fileNameFS)
	assert.NoError(t, err)

	err = filesystem.Remove(fileNameFS)
	assert.NoError(t, err)

	_, err = filesystem.Stat(fileName)
	assert.NotNil(t, err)
}

func TestCloseFile(t *testing.T) {
	testFile, _ := filesystem.Create(fileNameFS)
	filesystem.CloseFile(testFile)
	err := testFile.Close()

	assert.NotNil(t, err)
}

func TestWriteToWriter(t *testing.T) {
	content := []byte("{}")
	testFile, _ := filesystem.Create(fileNameFS)
	defer deleteFile(t, testFile)

	_, err := filesystem.WriteToWriter(testFile, content)

	assert.NoError(t, err)
	fileContents, err := os.ReadFile(fileNameFS)
	assert.NoError(t, err)
	assert.Equal(t, fileContents, content)
}
