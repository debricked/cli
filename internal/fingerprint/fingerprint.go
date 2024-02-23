package fingerprint

import (
	"archive/zip"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/tui"
	"lukechampine.com/blake3"
)

var EXCLUDED_EXT = []string{
	".1", ".2", ".3", ".4", ".5", ".6", ".7", ".8", ".9", ".ac", ".adoc", ".am",
	".asciidoc", ".bmp", ".build", ".cfg", ".chm", ".cmake", ".cnf",
	".conf", ".config", ".contributors", ".copying", ".crt", ".csproj", ".css",
	".csv", ".dat", ".data", ".doc", ".docx", ".dtd", ".dts", ".iws", ".c9", ".c9revisions",
	".dtsi", ".dump", ".eot", ".eps", ".geojson", ".gdoc", ".gif",
	".glif", ".gmo", ".gradle", ".guess", ".hex", ".htm", ".html", ".ico", ".iml",
	".in", ".inc", ".info", ".ini", ".ipynb", ".jpeg", ".jpg", ".json", ".jsonld", ".lock",
	".log", ".m4", ".map", ".markdown", ".md", ".md5", ".meta", ".mk", ".mxml",
	".o", ".otf", ".out", ".pbtxt", ".pdf", ".pem", ".phtml", ".plist", ".png",
	".po", ".ppt", ".prefs", ".properties", ".pyc", ".qdoc", ".result", ".rgb",
	".rst", ".scss", ".sha", ".sha1", ".sha2", ".sha256", ".sln", ".spec", ".sql",
	".sub", ".svg", ".svn-base", ".tab", ".template", ".test", ".tex", ".tiff",
	".toml", ".ttf", ".txt", ".utf-8", ".vim", ".wav", ".whl", ".woff", ".woff2", ".xht",
	".xhtml", ".xls", ".xlsx", ".xpm", ".xsd", ".xul", ".yaml", ".yml", ".wfp",
	".editorconfig", ".dotcover", ".pid", ".lcov", ".egg", ".manifest", ".cache", ".coverage", ".cover",
	".gem", ".lst", ".pickle", ".pdb", ".gml", ".pot", ".plt",
}

var EXCLUDED_FILE_ENDINGS = []string{"-doc", "changelog", "config", "copying", "license", "authors", "news", "licenses", "notice",
	"readme", "swiftdoc", "texidoc", "todo", "version", "ignore", "manifest", "sqlite", "sqlite3"}

var EXCLUDED_FILES = []string{
	"gradlew", "gradlew.bat", "mvnw", "mvnw.cmd", "gradle-wrapper.jar", "maven-wrapper.jar",
	"thumbs.db", "babel.config.js", "license.txt", "license.md", "copying.lib", "makefile",
	"[content_types].xml",
}

var FILES_TO_UNPACK = []string{".jar", ".nupkg", ".war"}

const HASH_SIZE = 16

func newHasher() *blake3.Hasher {
	return blake3.New(
		HASH_SIZE,
		nil,
	)
}

const (
	OutputFileNameFingerprints = "debricked.fingerprints.txt"
)

func isExcludedFile(filename string) bool {

	return isExcludedByExtension(filename) ||
		isExcludedByFilename(filename) ||
		isExcludedByEnding(filename)
}

func isExcludedByExtension(filename string) bool {
	filenameLower := strings.ToLower(filename)
	for _, format := range EXCLUDED_EXT {
		if filepath.Ext(filenameLower) == format {
			return true
		}
	}

	return false
}

func isExcludedByFilename(filename string) bool {
	filenameLower := strings.ToLower(filename)
	for _, file := range EXCLUDED_FILES {
		if filenameLower == file {
			return true
		}
	}

	return false
}

func isExcludedByEnding(filename string) bool {
	filenameLower := strings.ToLower(filename)
	for _, ending := range EXCLUDED_FILE_ENDINGS {
		if strings.HasSuffix(filenameLower, ending) {
			return true
		}
	}

	return false
}

type IFingerprint interface {
	FingerprintFiles(rootPath string, exclusions []string, fingerprintCompressedContent bool, includeSnippet bool) (Fingerprints, error)
}

type Fingerprinter struct {
	spinnerManager tui.ISpinnerManager
}

func NewFingerprinter() *Fingerprinter {
	return &Fingerprinter{
		spinnerManager: tui.NewSpinnerManager("Fingerprinting", "0"),
	}
}

type FileFingerprint struct {
	path          string
	contentLength int64
	fingerprint   []byte
	snippets      []Snippet
}

func (f FileFingerprint) ToString() string {
	path := filepath.ToSlash(f.path)
	output := fmt.Sprintf("file=%x,%d,%s", f.fingerprint, f.contentLength, path)

	if len(f.snippets) != 0 {
		snippetOutput := ""
		prevLine := -1
		for _, snippet := range f.snippets {
			if snippet.Line == prevLine {
				snippetOutput += fmt.Sprintf(",%s", snippet.Hash)
			} else {
				snippetOutput += fmt.Sprintf("\n%d=%s", snippet.Line, snippet.Hash)
			}
			prevLine = snippet.Line
		}

		output += snippetOutput
	}

	return output
}

func (f *Fingerprinter) FingerprintFiles(rootPath string, exclusions []string, fingerprintCompressedContent bool, includeSnippet bool) (Fingerprints, error) {
	if len(rootPath) == 0 {
		rootPath = filepath.Base("")
	}

	fingerprints := Fingerprints{}

	f.spinnerManager.Start()
	spinnerMessage := "files processed"
	spinner := f.spinnerManager.AddSpinner(spinnerMessage)

	nbFiles := 0

	err := filepath.Walk(rootPath, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fingerprintsZip, err := computeHashForFileAndZip(fileInfo, path, exclusions, fingerprintCompressedContent, includeSnippet)
		if err != nil {
			return err
		}
		if len(fingerprintsZip) != 0 {
			fingerprints.Entries = append(fingerprints.Entries, fingerprintsZip...)

			nbFiles += len(fingerprintsZip)

			if nbFiles%100 == 0 {
				f.spinnerManager.SetSpinnerMessage(spinner, spinnerMessage, fmt.Sprintf("%d", nbFiles))
			}
		}

		return nil
	})

	f.spinnerManager.SetSpinnerMessage(spinner, spinnerMessage, fmt.Sprintf("%d", nbFiles))

	if err != nil {
		spinner.Error()
	} else {
		spinner.Complete()
	}

	f.spinnerManager.Stop()

	return fingerprints, err
}

func computeHashForFileAndZip(fileInfo os.FileInfo, path string, exclusions []string, fingerprintCompressedContent bool, includeSnippet bool) ([]FileFingerprint, error) {
	if !shouldProcessFile(fileInfo, exclusions, path) {
		return nil, nil
	}

	var fingerprints []FileFingerprint

	// If the file should be unzipped, try to unzip and fingerprint it
	if isCompressedFile(path) && fingerprintCompressedContent {
		fingerprintsZip, err := inMemFingerprintingCompressedContent(path, exclusions, includeSnippet)
		if err != nil {
			if errors.Is(err, zip.ErrFormat) {
				fmt.Printf("WARNING: Could not unpack and fingerprint contents of compressed file [%s]. Error: %v\n", path, err)
			} else {
				return nil, err
			}
		}
		fingerprints = append(fingerprints, fingerprintsZip...)
	}

	fingerprint, err := computeHashForFile(path, includeSnippet)
	if err != nil {
		return nil, err
	}

	return append(fingerprints, fingerprint), nil
}

func isSymlink(filename string) (bool, error) {
	info, err := os.Lstat(filename)
	if err != nil {
		return false, err
	}

	return info.Mode()&os.ModeSymlink != 0, nil
}

var isSymlinkFunc = isSymlink

func shouldProcessFile(fileInfo os.FileInfo, exclusions []string, path string) bool {
	if fileInfo.IsDir() {
		return false
	}

	if file.Excluded(exclusions, path) {
		return false
	}

	if isExcludedFile(path) {
		return false
	}

	isSymlink, err := isSymlinkFunc(path)
	if err != nil {
		// Handle error with reading inmem files in windows
		if strings.HasSuffix(err.Error(), "The system cannot find the path specified.") {
			return true
		}
		// If we get a "not a directory" error, we can assume it's not a symlink
		// otherwise, we don't know, so we return false
		return strings.HasSuffix(err.Error(), "not a directory")
	}

	return !isSymlink
}

func computeHashForFile(filename string, includeSnippet bool) (FileFingerprint, error) {
	rc, err := os.Open(filename)
	if err != nil {
		return FileFingerprint{}, err
	}
	defer rc.Close()

	hasher := newHasher()

	contentLen, err := io.Copy(hasher, rc) // #nosec
	if err != nil {
		return FileFingerprint{}, err
	}

	if err != nil {
		return FileFingerprint{}, err
	}

	fingerprint := FileFingerprint{
		path:          filename,
		contentLength: contentLen,
		fingerprint:   hasher.Sum(nil),
	}

	if includeSnippet {
		winnowing := NewWinnowing(nil)

		snippets, err := winnowing.GenerateWFP(filename)
		if err != nil {
			return FileFingerprint{}, err
		}

		if snippets != nil {
			fingerprint.snippets = *snippets
		}
	}

	return fingerprint, nil
}

type Fingerprints struct {
	Entries []FileFingerprint `json:"fingerprints"`
}

func (f *Fingerprints) Len() int {
	return len(f.Entries)
}

var osCreate = os.Create

func (f *Fingerprints) ToFile(ouputFile string) error {
	file, err := osCreate(ouputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, fingerprint := range f.Entries {
		_, err := writer.WriteString(fingerprint.ToString() + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()

	return nil

}

func isCompressedFile(filename string) bool {
	for _, file := range FILES_TO_UNPACK {
		if filepath.Ext(filename) == file {
			return true
		}
	}

	return false
}

type NotProcessableError struct {
	Message string
}

func (e *NotProcessableError) Error() string {
	return e.Message
}

func processFile(f *zip.File, filename string, exclusions []string, includeSnippet bool) (*FileFingerprint, error) {
	longFileName := filepath.Join(filename, f.Name) // #nosec

	if !shouldProcessFile(f.FileInfo(), exclusions, longFileName) {

		return nil, &NotProcessableError{
			Message: "file is not processable",
		}
	}
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}

	defer rc.Close()

	hasher := newHasher()

	contentLen, err := io.Copy(hasher, rc) // #nosec
	if err != nil {
		return nil, err
	}

	fingerprint := FileFingerprint{
		path:          longFileName,
		contentLength: contentLen,
		fingerprint:   hasher.Sum(nil),
	}
	if includeSnippet {
		winnowing := NewWinnowing(nil)

		snippets, err := winnowing.GenerateWFP(longFileName)
		if err != nil {
			return nil, err
		}

		if snippets != nil {
			fingerprint.snippets = *snippets
		}
	}

	return &fingerprint, nil
}

func inMemFingerprintingCompressedContent(filename string, exclusions []string, includeSnippet bool) ([]FileFingerprint, error) {

	r, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	fingerprints := []FileFingerprint{}

	for _, f := range r.File {
		if filepath.IsAbs(f.Name) || strings.HasPrefix(f.Name, "..") {
			continue
		}

		fingerprint, err := processFile(f, filename, exclusions, includeSnippet)
		if err != nil {
			if _, ok := err.(*NotProcessableError); ok {
				continue
			} else {
				return nil, err
			}
		}
		fingerprints = append(fingerprints, *fingerprint)
	}

	return fingerprints, nil
}
