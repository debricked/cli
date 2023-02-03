package writer

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var fw = FileWriter{}

const fileName = "debricked-test.json"

func TestCreate(t *testing.T) {
	defer deleteFile(t, fileName)

	testFile, err := fw.Create(fileName)

	assert.NoError(t, err)
	assert.NotNil(t, testFile)
}

func TestWrite(t *testing.T) {
	defer deleteFile(t, fileName)
	testFile, _ := fw.Create(fileName)
	content := []byte("{}")

	err := fw.Write(testFile, content)

	assert.NoError(t, err)
	fileContents, err := os.ReadFile(fileName)
	assert.NoError(t, err)
	assert.Equal(t, fileContents, content)
}

func TestClose(t *testing.T) {
	defer deleteFile(t, fileName)
	testFile, _ := fw.Create(fileName)

	err := fw.Close(testFile)

	assert.NoError(t, err)
}

func deleteFile(t *testing.T, name string) {
	err := os.Remove(name)
	assert.NoError(t, err)
}
