package fingerprint

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

var errorString = "mock error"

// Test errors in symlink
func mockSymlink(filename string) (bool, error) {
	return false, fmt.Errorf(errorString)
}
func TestShouldProcessFile(t *testing.T) {
	// Create a temporary directory to use for testing
	tempDir, err := os.MkdirTemp("", "should-process-file-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test file and a symbolic link to the file in the temporary directory
	testFile := filepath.Join(tempDir, "test.py")
	if err := os.WriteFile(testFile, []byte("test"), 0600); err != nil {
		t.Fatalf("Failed to create test file %s: %v", testFile, err)
	}
	testLink := filepath.Join(tempDir, "test-link.py")
	if err := os.Symlink(testFile, testLink); err != nil {
		t.Fatalf("Failed to create symbolic link %s: %v", testLink, err)
	}

	tests := []struct {
		name     string
		filePath string
		excludes []string
		mock     func()
		want     bool
	}{
		{
			name:     "Test with a regular file",
			filePath: testFile,
			excludes: []string{},
			mock:     func() {},
			want:     true,
		},
		{
			name:     "Test with a symbolic link",
			filePath: testLink,
			excludes: []string{},
			mock:     func() {},
			want:     false,
		},
		{
			name:     "Test Excluded",
			filePath: testFile,
			excludes: []string{"**/test.py"},
			mock:     func() {},
			want:     false,
		},
		{
			name:     "Test with mockSymlink",
			filePath: testFile,
			excludes: []string{},
			mock:     func() { isSymlinkFunc = mockSymlink },
			want:     false,
		},
		{
			name:     "Test with errorString: The system cannot find the path specified.",
			filePath: testFile,
			excludes: []string{},
			mock:     func() { errorString = "The system cannot find the path specified." },
			want:     true,
		},
		{
			name:     "Test with errorString: not a directory",
			filePath: testFile,
			excludes: []string{},
			mock:     func() { errorString = "not a directory" },
			want:     true,
		},
		{
			name:     "Test with generic error",
			filePath: testFile,
			excludes: []string{},
			mock:     func() { errorString = "generic error" },
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() { errorString = "mock error" }()
			tt.mock()
			fileInfo, err := os.Stat(tt.filePath)
			if err != nil {
				t.Fatalf("Failed to get file info for %s: %v", tt.filePath, err)
			}
			if got := shouldProcessFile(fileInfo, tt.excludes, tt.filePath); got != tt.want {
				t.Errorf("Expected shouldProcessFile to return %v for %s, but it returned %v", tt.want, tt.filePath, got)
			}
		})
	}

	// Reset isSymlinkFunc and errorString
	isSymlinkFunc = isSymlink
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

	tempDir, err := os.MkdirTemp("", "slash-test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	testFile := filepath.Join(tempDir, "testfile.py")

	testFileSlashes := filepath.ToSlash(testFile)

	fingerprint := FileFingerprint{
		path:          testFile,
		contentLength: 21,
		fingerprint:   []byte{114, 33, 77, 180, 225, 229, 67, 1, 141, 27, 175, 232, 110, 163, 180, 68, 68, 68, 68, 68, 68},
	}

	assert.Equal(t, fmt.Sprintf("file=72214db4e1e543018d1bafe86ea3b4444444444444,21,%s", testFileSlashes), fingerprint.ToString())

	// Make sure it only contains "/" and not "\"
	assert.NotContains(t, fingerprint.ToString(), "\\")
	assert.Contains(t, fingerprint.ToString(), "/")

}

func TestFileFingerprintToString(t *testing.T) {
	fileFingerprint := FileFingerprint{path: "path", contentLength: 10, fingerprint: []byte("fingerprint")}
	assert.Equal(t, "file=66696e6765727072696e74,10,path", fileFingerprint.ToString())
}

func TestComputeMD5(t *testing.T) {
	// Test file not found
	_, err := computeMD5ForFile("testdata/fingerprinter/testfile-not-found.py")
	assert.Error(t, err)

	// Test file found
	entry, err := computeMD5ForFile("testdata/fingerprinter/testfile.py")
	assert.NoError(t, err)
	entryS := fmt.Sprintf("%x", entry.fingerprint)
	assert.Equal(t, "72214db4e1e543018d1bafe86ea3b444", entryS)
}

func TestFingerprintsToFile(t *testing.T) {
	tests := []struct {
		name          string
		outputFile    string
		setupMock     func()
		expectedError bool
	}{
		{
			name:          "Successful write",
			outputFile:    "fingerprints.wfp",
			setupMock:     func() {},
			expectedError: false,
		},
		{
			name: "Failed to create file",
			setupMock: func() {
				osCreate = func(name string) (*os.File, error) {
					return nil, errors.New("forced error")
				}
			},
			outputFile:    "test/fingerprints.wfp",
			expectedError: true,
		},
		{
			name: "Failed to write to file",
			setupMock: func() {
				osCreate = func(name string) (*os.File, error) {
					return os.Create("test/fingerprints.wfp")
				}
			},
			outputFile:    "/invalid/path/fingerprints.wfp",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset osCreate to its original function after each test
			defer func() { osCreate = os.Create }()

			// Setup the mock function
			tt.setupMock()

			// Create temp dir
			dir, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Fatalf("Failed to create temporary directory: %v", err)
			}
			defer os.RemoveAll(dir)

			// Create fingerprints
			fingerprints := Fingerprints{}
			fingerprints.Entries = append(fingerprints.Entries, FileFingerprint{path: "path", contentLength: 10, fingerprint: []byte("fingerprint")})

			// Write fingerprints to file
			err = fingerprints.ToFile(filepath.Join(dir, tt.outputFile))
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestShouldUnzip(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "Should unzip .jar file",
			filename: "test.jar",
			want:     true,
		},
		{
			name:     "Should unzip .nupkg file",
			filename: "test.nupkg",
			want:     true,
		},
		{
			name:     "Should not unzip .txt file",
			filename: "test.txt",
			want:     false,
		},
		{
			name:     "Should not unzip .go file",
			filename: "test.go",
			want:     false,
		},
		{
			name:     "Should pick up .jar file in nested folder",
			filename: "deep/folder/test.jar",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldUnzip(tt.filename); got != tt.want {
				t.Errorf("shouldUnzip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemFingerprintingCompressedContent(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected int
		suffix   string
	}{
		{
			name:     "Jar",
			path:     "testdata/zipfile/jar",
			expected: 5,
			suffix:   "log4j-api-2.18.0.jar",
		},
		{
			name:     "Nupkg",
			path:     "testdata/zipfile/nupkg",
			expected: 22,
			suffix:   "newtonsoft.json.13.0.3.nupkg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fingerprinter := NewFingerprinter()
			fingerprints, err := fingerprinter.FingerprintFiles(tt.path, []string{})
			assert.NoError(t, err)
			assert.NotNil(t, fingerprints)
			assert.NotEmpty(t, fingerprints)
			assert.Equal(t, tt.expected, fingerprints.Len())
			lastRow := fingerprints.Entries[len(fingerprints.Entries)-1]
			assert.True(t, strings.HasSuffix(lastRow.ToString(), tt.suffix))
		})
	}
}

func TestComputeMD5ForFile(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "Non-existent file",
			file:    "non_existent_file.txt",
			wantErr: true,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := computeMD5ForFile(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("computeMD5ForFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
