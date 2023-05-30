package io

import (
	"embed"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var filesystem = FileSystem{}

const fileNameFS = "debricked-test-fs"

//go:embed testdata/embed-file
var embedFile embed.FS

const embedFilePath = "testdata/embed-file"

func TestCreateFile(t *testing.T) {
	fn := fileNameFS + t.Name()
	testFile, err := filesystem.Create(fn)
	assert.NoError(t, err)
	assert.NotNil(t, testFile)
	defer deleteFile(t, testFile)
}

func TestOpen(t *testing.T) {
	fn := fileNameFS + t.Name()
	testFile, _ := filesystem.Create(fn)
	defer deleteFile(t, testFile)

	_, err := filesystem.Open(fn)

	assert.NoError(t, err)
}

func TestStat(t *testing.T) {
	fn := fileNameFS + t.Name()
	_, err := filesystem.Stat(fn)
	assert.NotNil(t, err)

	testFile, _ := filesystem.Create(fn)
	defer deleteFile(t, testFile)

	_, err = filesystem.Stat(fn)

	assert.NoError(t, err)
}

func TestStatFile(t *testing.T) {
	fn := fileNameFS + t.Name()
	testFile, _ := filesystem.Create(fn)
	defer deleteFile(t, testFile)

	_, err := filesystem.StatFile(testFile)

	assert.NoError(t, err)
}

func TestReadFile(t *testing.T) {
	fn := fileNameFS + t.Name()
	testFile, _ := filesystem.Create(fn)
	defer deleteFile(t, testFile)

	_, err := filesystem.ReadFile(fn)

	assert.NoError(t, err)
}

func TestRemoveFile(t *testing.T) {
	fn := fileNameFS + t.Name()
	testFile, _ := filesystem.Create(fn)
	defer testFile.Close()

	_, err := filesystem.Stat(fn)
	assert.NoError(t, err)

	err = filesystem.Remove(fn)
	assert.NoError(t, err)

	_, err = filesystem.Stat(fn)
	assert.NotNil(t, err)
}

func TestCloseFile(t *testing.T) {
	fn := fileNameFS + t.Name()
	testFile, _ := filesystem.Create(fn)
	defer deleteFile(t, testFile)
	filesystem.CloseFile(testFile)
	err := testFile.Close()

	assert.NotNil(t, err)
}

func TestWriteToWriter(t *testing.T) {
	fn := fileNameFS + t.Name()
	content := []byte("{}")
	testFile, _ := filesystem.Create(fn)
	defer deleteFile(t, testFile)

	_, err := filesystem.WriteToWriter(testFile, content)

	assert.NoError(t, err)
	fileContents, err := os.ReadFile(fn)
	assert.NoError(t, err)
	assert.Equal(t, fileContents, content)
}

func TestMkdirTemp(t *testing.T) {
	fn := fileNameFS + t.Name()
	tmpdir, err := filesystem.MkdirTemp(fn)
	defer os.RemoveAll(tmpdir)
	assert.NoError(t, err)
	assert.NotNil(t, tmpdir)
}

func TestRemoveAll(t *testing.T) {
	fn := fileNameFS + t.Name()
	tmpdir, _ := filesystem.MkdirTemp(fn)
	_, err := filesystem.Stat(tmpdir)
	assert.Nil(t, err)
	filesystem.RemoveAll(tmpdir)
	_, err = filesystem.Stat(tmpdir)
	assert.NotNil(t, err)
}

func TestOpenEmbed(t *testing.T) {
	file, err := filesystem.FsOpenEmbed(embedFile, embedFilePath)
	assert.Nil(t, err)
	defer file.Close()
}

func TestCloseFs(t *testing.T) {
	file, _ := filesystem.FsOpenEmbed(embedFile, embedFilePath)
	filesystem.FsCloseFile(file)
}

func TestReadAll(t *testing.T) {
	file, _ := filesystem.FsOpenEmbed(embedFile, embedFilePath)
	defer file.Close()
	bytes, err := filesystem.FsReadAll(file)

	assert.Nil(t, err)
	assert.NotNil(t, bytes)
}

func TestWriteFile(t *testing.T) {
	fn := fileNameFS + t.Name()
	err := filesystem.FsWriteFile(fn, []byte{}, 0600)
	_ = filesystem.Remove(fn)

	assert.Nil(t, err)

}
