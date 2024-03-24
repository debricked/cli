package fingerprint

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/bzip2"
	"compress/gzip"
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

type DebrickedOptions struct {
	Path                         string
	Exclusions                   []string
	Inclusions                   []string
	FingerprintCompressedContent bool
	MinFingerprintContentLength  int
}

var ZIP_FILE_ENDINGS = []string{".jar", ".nupkg", ".war", ".zip", ".ear", ".whl"}
var TAR_GZIP_FILE_ENDINGS = []string{".tgz", ".tar.gz"}
var TAR_BZIP2_FILE_ENDINGS = []string{".tar.bz2"}

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

type IFingerprint interface {
	FingerprintFiles(
		options DebrickedOptions,
	) (Fingerprints, error)
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

func (f *Fingerprinter) FingerprintFiles(options DebrickedOptions) (Fingerprints, error) {
	if len(options.Path) == 0 {
		options.Path = filepath.Base("")
	}

	fingerprints := Fingerprints{}

	f.spinnerManager.Start()
	spinnerMessage := "files processed"
	spinner := f.spinnerManager.AddSpinner(spinnerMessage)

	nbFiles := 0

	err := filepath.Walk(options.Path, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		fileFingerprints, err := computeHashForFileAndZip(fileInfo, path, options.Exclusions, options.Inclusions, options.FingerprintCompressedContent)
		if err != nil {
			return err
		}

		var filteredFileFingerprints []FileFingerprint
		for _, fileFingerprint := range fileFingerprints {
			if fileFingerprint.contentLength >= int64(options.MinFingerprintContentLength) {
				filteredFileFingerprints = append(filteredFileFingerprints, fileFingerprint)
			}
		}

		if len(filteredFileFingerprints) != 0 {
			fingerprints.Entries = append(fingerprints.Entries, filteredFileFingerprints...)

			nbFiles += len(filteredFileFingerprints)

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

func computeHashForArchive(path string, exclusions []string, inclusions []string) ([]FileFingerprint, error) {
	if isZipFile(path) {
		return inMemFingerprintZipContent(path, exclusions, inclusions)
	}
	if isTarGZipFile(path) {
		return inMemFingerprintTarGZipContent(path, exclusions, inclusions)
	}
	if isTarBZip2File(path) {
		return inMemFingerprintTarBZip2Content(path, exclusions, inclusions)
	}

	return nil, nil
}

func computeHashForFileAndZip(
	fileInfo os.FileInfo, path string, exclusions []string, inclusions []string, fingerprintCompressedContent bool,
) ([]FileFingerprint, error) {
	if !shouldProcessFile(fileInfo, exclusions, inclusions, path) {
		return nil, nil
	}

	var fingerprints []FileFingerprint

	if fingerprintCompressedContent {
		fingerprintsArchive, err := computeHashForArchive(path, exclusions, inclusions)
		if err != nil {
			if errors.Is(err, zip.ErrFormat) {
				fmt.Printf("WARNING: Could not unpack and fingerprint contents of compressed file [%s]. Error: %v\n", path, err)
			} else {
				return nil, err
			}
		}
		fingerprints = append(fingerprints, fingerprintsArchive...)
	}

	fingerprint, err := computeHashForFile(path)
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

func shouldProcessFile(fileInfo os.FileInfo, exclusions []string, inclusions []string, path string) bool {
	inclusions = append(inclusions, DefaultInclusionsFingerprint()...)
	exclusions = append(exclusions, DefaultExclusionsFingerprint()...)
	if fileInfo.IsDir() {
		return false
	}
	if file.Excluded(exclusions, inclusions, path) {
		return false
	}
	if !strings.Contains(filepath.Base(path), ".") {
		return false
	} // If no extension

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

func computeHashForFile(filename string) (FileFingerprint, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return FileFingerprint{}, err
	}

	hasher := newHasher()

	if _, err := hasher.Write(data); err != nil {
		return FileFingerprint{}, err
	}

	contentLength := int64(len(data))

	return FileFingerprint{
		path:          filename,
		contentLength: contentLength,
		fingerprint:   hasher.Sum(nil),
	}, nil
}

type Fingerprints struct {
	Entries []FileFingerprint `json:"fingerprints"`
}

func (f *Fingerprints) Len() int {
	return len(f.Entries)
}

var osCreate = os.Create

func ensureDirExists(dir string) error {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		return nil
	}

	return os.MkdirAll(dir, 0755)
}

func (f *Fingerprints) ToFile(outputFile string) error {
	dir := filepath.Dir(outputFile)
	if err := ensureDirExists(dir); err != nil {
		return fmt.Errorf("failed to ensure directory exists: %w", err)
	}

	file, err := osCreate(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	return f.writeToFile(file)
}

func (f *Fingerprints) writeToFile(file *os.File) error {
	writer := bufio.NewWriter(file)
	for _, fingerprint := range f.Entries {
		if _, err := writer.WriteString(fingerprint.ToString() + "\n"); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}

	return writer.Flush()
}
func isZipFile(filename string) bool {
	for _, file := range ZIP_FILE_ENDINGS {
		if filepath.Ext(filename) == file {
			return true
		}
	}

	return false
}

func isTarGZipFile(filename string) bool {
	for _, file := range TAR_GZIP_FILE_ENDINGS {
		if strings.HasSuffix(filename, file) {
			return true
		}
	}

	return false
}

func isTarBZip2File(filename string) bool {
	for _, file := range TAR_BZIP2_FILE_ENDINGS {
		if strings.HasSuffix(filename, file) {
			return true
		}
	}

	return false
}

func shouldProcessTarHeader(header tar.Header, exclusions []string, inclusions []string, longPath string) bool {
	if header.Typeflag != tar.TypeReg {
		return false
	}
	if filepath.IsAbs(header.Name) || strings.HasPrefix(header.Name, "..") {
		return false
	}
	if !shouldProcessFile(header.FileInfo(), exclusions, inclusions, longPath) {
		return false
	}

	return true
}

func inMemFingerprintTarBZip2Content(filename string, exclusions []string, inclusions []string) ([]FileFingerprint, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bz2Reader := bzip2.NewReader(file)
	tarReader := tar.NewReader(bz2Reader)
	fingerprints := []FileFingerprint{}
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		longPath := filepath.Join(filename, header.Name) // #nosec
		fmt.Println("Extracted:", longPath)
		if !shouldProcessTarHeader(*header, exclusions, inclusions, longPath) {
			continue
		}
		hasher := newHasher()

		_, err = io.Copy(hasher, tarReader) // #nosec
		if err != nil {
			return nil, err
		}

		fingerprints = append(fingerprints, FileFingerprint{
			path:          longPath,
			contentLength: header.Size,
			fingerprint:   hasher.Sum(nil),
		})
	}

	return fingerprints, nil
}

func inMemFingerprintTarGZipContent(filename string, exclusions []string, inclusions []string) ([]FileFingerprint, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return nil, err
	}
	tarReader := tar.NewReader(gzReader)
	fingerprints := []FileFingerprint{}
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		longPath := filepath.Join(filename, header.Name) // #nosec
		if !shouldProcessTarHeader(*header, exclusions, inclusions, longPath) {
			continue
		}
		hasher := newHasher()

		_, err = io.Copy(hasher, tarReader) // #nosec
		if err != nil {
			return nil, err
		}

		fingerprints = append(fingerprints, FileFingerprint{
			path:          longPath,
			contentLength: header.Size,
			fingerprint:   hasher.Sum(nil),
		})
	}

	return fingerprints, nil
}

func inMemFingerprintZipContent(filename string, exclusions []string, inclusions []string) ([]FileFingerprint, error) {
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

		if !shouldProcessFile(f.FileInfo(), exclusions, inclusions, longFileName) {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}

		hasher := newHasher()

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
