package file

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsExcludedFile(t *testing.T) {

	// Test excluded file extensions
	excludedExts := []string{".doc", ".pdf", ".txt"}
	for _, ext := range excludedExts {
		filename := "file" + ext
		assert.True(t, isExcludedFile(filename), "Expected %q to be excluded", filename)
	}

	// Test excluded files
	excludedFiles := []string{"LICENSE", "README.md", "Makefile"}
	for _, file := range excludedFiles {
		assert.True(t, isExcludedFile(file), "Expected %q to be excluded", file)
	}

	// Test excluded file endings
	excludedEndings := []string{"-doc", "changelog", "config", "copying", "license", "authors", "news", "licenses", "notice",
		"readme", "swiftdoc", "texidoc", "todo", "version", "ignore", "manifest", "sqlite", "sqlite3"}
	for _, ending := range excludedEndings {
		filename := "file." + ending
		assert.True(t, isExcludedFile(filename), "Expected %q to be excluded", filename)
	}

	// Test non-excluded files
	assert.False(t, isExcludedFile("file.py"), "Expected file.txt to not be excluded")
	assert.False(t, isExcludedFile("file.go"), "Expected .go to not be excluded")
	assert.False(t, isExcludedFile("file.dll"), "Expected .dll to not be excluded")
	assert.False(t, isExcludedFile("file.jar"), "Expected .jar to not be excluded")
}

func TestNewFingerprinter(t *testing.T) {
	assert.NotNil(t, NewFingerprinter())
}

func TestFingerprinterInterface(t *testing.T) {
	assert.Implements(t, (*IFingerprint)(nil), new(Fingerprinter))
}

func TestFingerprintFiles(t *testing.T) {
	fingerprinter := NewFingerprinter()
	fingerprints, err := fingerprinter.FingerprintFiles("testdata/fingerprinter", []string{})
	assert.NoError(t, err)
	assert.NotNil(t, fingerprints)
	assert.NotEmpty(t, fingerprints)
	assert.Equal(t, 1, fingerprints.Len())
	assert.Equal(t, "file=72214db4e1e543018d1bafe86ea3b444,21,testdata/fingerprinter/testfile.py", fingerprints.Entries[0].ToString())

	// Test no file
	fingerprints, err = fingerprinter.FingerprintFiles("", []string{})
	assert.NoError(t, err)
	assert.NotNil(t, fingerprints)
	assert.NotEmpty(t, fingerprints)

}

func TestFingerprintFilesBackslash(t *testing.T) {
	fingerprint := FileFingerprint{
		path:          "testdata\\fingerprinter\\testfile.py",
		contentLength: 21,
		fingerprint:   []byte{114, 33, 77, 180, 225, 229, 67, 1, 141, 27, 175, 232, 110, 163, 180, 68, 68, 68, 68, 68, 68},
	}

	assert.Equal(t, "file=72214db4e1e543018d1bafe86ea3b4444444444444,21,testdata/fingerprinter/testfile.py", fingerprint.ToString())

}

func TestFileFingerprintToString(t *testing.T) {
	fileFingerprint := FileFingerprint{path: "path", contentLength: 10, fingerprint: []byte("fingerprint")}
	assert.Equal(t, "file=66696e6765727072696e74,10,path", fileFingerprint.ToString())
}

func TestComputeMD5(t *testing.T) {
	// Test file not found
	_, err := computeMD5("testdata/fingerprinter/testfile-not-found.py")
	assert.Error(t, err)

	// Test file found
	entry, err := computeMD5("testdata/fingerprinter/testfile.py")
	assert.NoError(t, err)
	entryS := fmt.Sprintf("%x", entry.fingerprint)
	assert.Equal(t, "72214db4e1e543018d1bafe86ea3b444", entryS)
}

func TestFingerprintsToFile(t *testing.T) {
	fingerprints := Fingerprints{}
	fingerprints.Entries = append(fingerprints.Entries, FileFingerprint{path: "path", contentLength: 10, fingerprint: []byte("fingerprint")})
	// Create temp dir
	dir, err := os.MkdirTemp("", "test")
	assert.NoError(t, err)
	defer os.RemoveAll(dir)
	// Write fingerprints to file
	err = fingerprints.ToFile(dir + "/fingerprints.wfp")
	assert.NoError(t, err)

}
