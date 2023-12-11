package fingerprint

import (
	"archive/zip"
	"bufio"
	"crypto/md5" // #nosec
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/tui"
)

var EXCLUDED_EXT = []string{
	".1", ".2", ".3", ".4", ".5", ".6", ".7", ".8", ".9", ".ac", ".adoc", ".am",
	".asciidoc", ".bmp", ".build", ".cfg", ".chm", ".class", ".cmake", ".cnf",
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
}

const (
	OutputFileNameFingerprints = "debricked.fingerprints.wfp"
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
	FingerprintFiles(rootPath string, exclusions []string) (Fingerprints, error)
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
}

func (f FileFingerprint) ToString() string {
	path := filepath.ToSlash(f.path)

	return fmt.Sprintf("file=%x,%d,%s", f.fingerprint, f.contentLength, path)
}

func (f *Fingerprinter) FingerprintFiles(rootPath string, exclusions []string) (Fingerprints, error) {
	log.Println("Warning: Fingerprinting is beta and may not work as expected.")
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

		fingerprintsZip, err := computeMD5ForFileAndZip(fileInfo, path, exclusions)
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

func computeMD5ForFileAndZip(fileInfo os.FileInfo, path string, exclusions []string) ([]FileFingerprint, error) {
	if !shouldProcessFile(fileInfo, exclusions, path) {
		return nil, nil
	}

	var fingerprints []FileFingerprint

	// If the file should be unzipped, try to unzip and fingerprint it
	if shouldUnzip(path) {
		fingerprintsZip, err := inMemFingerprintingCompressedContent(path, exclusions)
		if err != nil {
			if strings.Contains(err.Error(), "not a valid zip file") {
				fmt.Printf("WARNING: Could not unpack and fingerprint contents of compressed file [%s]. Error: %v\n", path, err)
			} else {
				return nil, err
			}
		}
		fingerprints = append(fingerprints, fingerprintsZip...)
	}

	// Compute the MD5 for the file
	fingerprint, err := computeMD5ForFile(path)
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

func computeMD5ForFile(filename string) (FileFingerprint, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return FileFingerprint{}, err
	}

	hash := md5.New() // #nosec

	if _, err := hash.Write(data); err != nil {
		return FileFingerprint{}, err
	}

	contentLength := int64(len(data))

	if err != nil {
		return FileFingerprint{}, err
	}

	return FileFingerprint{
		path:          filename,
		contentLength: contentLength,
		fingerprint:   hash.Sum(nil),
	}, nil
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

var filesToUnzip = []string{".jar", ".nupkg"}

func shouldUnzip(filename string) bool {
	for _, file := range filesToUnzip {
		if filepath.Ext(filename) == file {
			return true
		}
	}

	return false
}

func inMemFingerprintingCompressedContent(filename string, exclusions []string) ([]FileFingerprint, error) {

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
		longFileName := filepath.Join(filename, f.Name) // #nosec

		if !shouldProcessFile(f.FileInfo(), exclusions, longFileName) {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		hasher := md5.New()          // #nosec
		_, err = io.Copy(hasher, rc) // #nosec
		if err != nil {
			rc.Close()

			return nil, err
		}

		fingerprints = append(fingerprints, FileFingerprint{
			path:          longFileName,
			contentLength: int64(f.UncompressedSize64),
			fingerprint:   hasher.Sum(nil),
		})

		rc.Close()
	}

	return fingerprints, nil
}
