package io

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fw = FileWriter{}

const fileName = "debricked-test.json"

func TestCreate(t *testing.T) {
	fn := fileName + t.Name()
	testFile, err := fw.Create(fn)
	assert.NoError(t, err)
	assert.NotNil(t, testFile)
	defer deleteFile(t, testFile)
}

func TestWrite(t *testing.T) {
	fn := fileName + t.Name()
	content := []byte("{}")
	testFile, _ := fw.Create(fn)
	defer deleteFile(t, testFile)

	err := fw.Write(testFile, content)

	assert.NoError(t, err)
	fileContents, err := os.ReadFile(fn)
	assert.NoError(t, err)
	assert.Equal(t, fileContents, content)
}

func TestClose(t *testing.T) {
	fn := fileName + t.Name()
	testFile, _ := fw.Create(fn)
	defer deleteFile(t, testFile)

	err := fw.Close(testFile)

	assert.NoError(t, err)
}

func deleteFile(t *testing.T, file *os.File) {
	_ = file.Close()
	err := os.Remove(file.Name())
	assert.NoError(t, err)
}
