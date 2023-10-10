package file

import (
	"bufio"
	"crypto/md5" // #nosec
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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
	".xhtml", ".xls", ".xlsx", ".xml", ".xpm", ".xsd", ".xul", ".yaml", ".yml", ".wfp",
	".editorconfig", ".dotcover", ".pid", ".lcov", ".egg", ".manifest", ".cache", ".coverage", ".cover",
	".gem", ".lst", ".pickle", ".pdb", ".gml", ".pot", ".plt",
}

var EXCLUDED_FILE_ENDINGS = []string{"-doc", "changelog", "config", "copying", "license", "authors", "news", "licenses", "notice",
	"readme", "swiftdoc", "texidoc", "todo", "version", "ignore", "manifest", "sqlite", "sqlite3"}

var ECLUDED_FILES = []string{
	"gradlew", "gradlew.bat", "mvnw", "mvnw.cmd", "gradle-wrapper.jar", "maven-wrapper.jar",
	"thumbs.db", "babel.config.js", "license.txt", "license.md", "copying.lib", "makefile",
}

const (
	OutputFileNameFingerprints = ".debricked.fingerprints.wfp"
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
	for _, file := range ECLUDED_FILES {
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
	// Replace backslashes with forward slashes to make the path platform independent
	path := strings.ReplaceAll(f.path, "\\", "/")

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
		nbFiles++

		if err != nil {
			return err
		}

		if !shouldProcessFile(fileInfo, exclusions, path) {
			return nil
		}

		fingerprint, err := computeMD5(path)
		if err != nil {
			return err
		}

		fingerprints.Append(fingerprint)

		if nbFiles%100 == 0 {
			f.spinnerManager.SetSpinnerMessage(spinner, spinnerMessage, fmt.Sprintf("%d", nbFiles))
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

func shouldProcessFile(fileInfo os.FileInfo, exclusions []string, path string) bool {
	if fileInfo.IsDir() {
		return false
	}

	if excluded(exclusions, path) {
		return false
	}

	if isExcludedFile(path) {
		return false
	}

	return true
}

func computeMD5(filename string) (FileFingerprint, error) {
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

func (f *Fingerprints) ToFile(ouputFile string) error {
	file, err := os.Create(ouputFile)
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

func (f *Fingerprints) Append(fingerprint FileFingerprint) {
	f.Entries = append(f.Entries, fingerprint)
}
