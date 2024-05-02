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
		includes []string
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
			includes: []string{},
			mock:     func() {},
			want:     false,
		},
		{
			name:     "Test Excluded",
			filePath: testFile,
			excludes: []string{"**/test.py"},
			includes: []string{},
			mock:     func() {},
			want:     false,
		},
		{
			name:     "Test Excluded and Included",
			filePath: testFile,
			excludes: []string{"**/test.py"},
			includes: []string{"**/test.py"},
			mock:     func() {},
			want:     true,
		},
		{
			name:     "Test with mockSymlink",
			filePath: testFile,
			excludes: []string{},
			includes: []string{},
			mock:     func() { isSymlinkFunc = mockSymlink },
			want:     false,
		},
		{
			name:     "Test with errorString: The system cannot find the path specified.",
			filePath: testFile,
			excludes: []string{},
			includes: []string{},
			mock:     func() { errorString = "The system cannot find the path specified." },
			want:     true,
		},
		{
			name:     "Test with errorString: not a directory",
			filePath: testFile,
			excludes: []string{},
			includes: []string{},
			mock:     func() { errorString = "not a directory" },
			want:     true,
		},
		{
			name:     "Test with generic error",
			filePath: testFile,
			excludes: []string{},
			includes: []string{},
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
			if got := shouldProcessFile(fileInfo, tt.excludes, tt.includes, tt.filePath); got != tt.want {
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
	fingerprints, err := fingerprinter.FingerprintFiles(
		DebrickedOptions{
			Path:                         "testdata/fingerprinter",
			Exclusions:                   []string{},
			Inclusions:                   []string{},
			FingerprintCompressedContent: true,
			MinFingerprintContentLength:  0,
		},
	)
	assert.NoError(t, err)
	assert.NotNil(t, fingerprints)
	assert.NotEmpty(t, fingerprints)
	assert.Equal(t, 2, fingerprints.Len())
	assert.Equal(t, "file=634c5485de8e22b27094affadd8a6e3b,21,testdata/fingerprinter/testfile.py", fingerprints.Entries[0].ToString())

	// Test no file
	fingerprints, err = fingerprinter.FingerprintFiles(
		DebrickedOptions{
			Path:                         "",
			Exclusions:                   []string{},
			Inclusions:                   []string{},
			FingerprintCompressedContent: true,
			MinFingerprintContentLength:  0,
		},
	)
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
	_, err := computeHashForFile("testdata/fingerprinter/testfile-not-found.py")
	assert.Error(t, err)

	// Test file found
	entry, err := computeHashForFile("testdata/fingerprinter/testfile.py")
	assert.NoError(t, err)
	entryS := fmt.Sprintf("%x", entry.fingerprint)
	assert.Equal(t, "634c5485de8e22b27094affadd8a6e3b", entryS)
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
		{
			name:          "Create non-existent directory",
			setupMock:     func() {},
			outputFile:    "test/newfile/debricked.fingerprints.txt",
			expectedError: false,
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

				// Check if the file exists
				_, err := os.Stat(filepath.Join(dir, tt.outputFile))
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
			name:     "Should unzip .whl file",
			filename: "test.whl",
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
			if got := isZipFile(tt.filename); got != tt.want {
				t.Errorf("isZipFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestIsTarGZip(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     ".jar is not tar gzip",
			filename: "test.jar",
			want:     false,
		},
		{
			name:     ".nupkg is not tar gzip",
			filename: "test.nupkg",
			want:     false,
		},
		{
			name:     ".txt is not tar gzip",
			filename: "test.txt",
			want:     false,
		},
		{
			name:     ".go is not tar gzip",
			filename: "test.go",
			want:     false,
		},
		{
			name:     "tar.bz2",
			filename: "deep/folder/python-dotenv-1.0.0.tar.bz2",
			want:     false,
		},
		{
			name:     ".tgz is tar gzip",
			filename: "test.tgz",
			want:     true,
		},
		{
			name:     "Should pick up .tgz archive in nested folder",
			filename: "deep/folder/test.tgz",
			want:     true,
		},
		{
			name:     "tar.gz",
			filename: "deep/folder/python-dotenv-1.0.0.tar.gz",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTarGZipFile(tt.filename); got != tt.want {
				t.Errorf("isTarGZipFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTarBZip2(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     ".zip is not tar bz2",
			filename: "test.zip",
			want:     false,
		},
		{
			name:     ".tgz is not tar bz2",
			filename: "test.tgz",
			want:     false,
		},
		{
			name:     ".tar.gz is not tar bz2",
			filename: "test.tar.gz",
			want:     false,
		},
		{
			name:     "Should pick up .tar.bz2 archive in nested folder",
			filename: "deep/folder/test.tar.bz2",
			want:     true,
		},
		{
			name:     "tar.bz2",
			filename: "deep/folder/python-dotenv-1.0.0.tar.bz2",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isTarBZip2File(tt.filename); got != tt.want {
				t.Errorf("isTarGZipFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInMemFingerprintingCompressedContent(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		expected    int
		suffix      string
		shouldUnzip bool
	}{
		{
			name:        "Jar",
			path:        "testdata/archive/jar",
			expected:    195,
			suffix:      "log4j-api-2.18.0.jar",
			shouldUnzip: true,
		},
		{
			name:        "Nupkg",
			path:        filepath.Join("testdata", "archive", "nupkg"),
			expected:    21,
			suffix:      "newtonsoft.json.13.0.3.nupkg",
			shouldUnzip: true,
		},
		{
			name:        "TGz",
			path:        "testdata/archive/tgz",
			expected:    984,
			suffix:      "lodash.tgz",
			shouldUnzip: true,
		},
		{
			name:        "Nupkg not unpack",
			path:        "testdata/archive/nupkg",
			expected:    1,
			suffix:      "newtonsoft.json.13.0.3.nupkg",
			shouldUnzip: false,
		},
		{
			name:        "BZip2",
			path:        "testdata/archive/bz2",
			expected:    7,
			suffix:      "stuf-0.1.tar.bz2",
			shouldUnzip: true,
		},
		{
			name:        "whl",
			path:        "testdata/archive/whl",
			expected:    19,
			suffix:      "requests-2.31.0-py3-none-any.whl",
			shouldUnzip: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fingerprinter := NewFingerprinter()
			fingerprints, err := fingerprinter.FingerprintFiles(
				DebrickedOptions{
					Path:                         tt.path,
					Exclusions:                   []string{},
					Inclusions:                   []string{},
					FingerprintCompressedContent: tt.shouldUnzip,
					MinFingerprintContentLength:  45,
				},
			)
			assert.NoError(t, err)
			assert.NotNil(t, fingerprints)
			assert.NotEmpty(t, fingerprints)
			assert.Equal(t, tt.expected, fingerprints.Len())
			lastRow := fingerprints.Entries[len(fingerprints.Entries)-1]
			assert.True(t, strings.HasSuffix(lastRow.ToString(), tt.suffix))
		})
	}
}

func TestComputeHashForFile(t *testing.T) {
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
			_, err := computeHashForFile(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("computeHashForFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
