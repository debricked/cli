package writer

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fw = FileWriter{}

const fileName = "debricked-test.json"

func TestCreate(t *testing.T) {
	testFile, err := fw.Create(fileName)
	assert.NoError(t, err)
	assert.NotNil(t, testFile)
	defer deleteFile(t, testFile)
}

func TestWrite(t *testing.T) {
	content := []byte("{}")
	testFile, _ := fw.Create(fileName)
	defer deleteFile(t, testFile)

	err := fw.Write(testFile, content)

	assert.NoError(t, err)
	fileContents, err := os.ReadFile(fileName)
	assert.NoError(t, err)
	assert.Equal(t, fileContents, content)
}

func TestClose(t *testing.T) {
	testFile, _ := fw.Create(fileName)
	defer deleteFile(t, testFile)

	err := fw.Close(testFile)

	assert.NoError(t, err)
}

func deleteFile(t *testing.T, file *os.File) {
	_ = file.Close()
	err := os.Remove(file.Name())
	assert.NoError(t, err)
}
