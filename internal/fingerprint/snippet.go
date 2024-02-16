package fingerprint

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const (
	Gram               = 30
	Window             = 64
	MinFileSize        = 256
	MaxPostSize        = 64 * 1024 // 64k
	MaxLongLineChars   = 1000
	ASCII0             = 48
	ASCII9             = 57
	ASCIILF            = 10
	ASCIIBackslash     = 92
	MaxCRC32           = 4294967295
	SkipSnippetExtSize = 29
)

var SkipSnippetExt = map[string]bool{
	".exe": true, ".zip": true, // Add all extensions as in the Python example
}

type Winnowing struct {
	sizeLimit      bool
	skipSnippets   bool
	maxPostSize    int
	allExtensions  bool
	obfuscate      bool
	fileMap        map[string]string
	crc8MaximTable []uint8
}

func NewWinnowing(sizeLimit, skipSnippets, allExtensions, obfuscate bool, postSize int) *Winnowing {
	return &Winnowing{
		sizeLimit:      sizeLimit,
		skipSnippets:   skipSnippets,
		maxPostSize:    postSize * 1024,
		allExtensions:  allExtensions,
		obfuscate:      obfuscate,
		fileMap:        make(map[string]string),
		crc8MaximTable: make([]uint8, 0),
	}
}

func (w *Winnowing) NormalizeByte(b byte) byte {
	if b < ASCII0 || b > ASCII9 {
		return 0
	}
	return b
}

func (w *Winnowing) ShouldSkipFile(filePath string) bool {
	extension := strings.ToLower(filepath.Ext(filePath))
	if _, ok := SkipSnippetExt[extension]; ok && !w.allExtensions {
		return true
	}
	return false
}

func (w *Winnowing) ReadFile(filePath string) ([]byte, error) {
	if w.ShouldSkipFile(filePath) {
		return nil, fmt.Errorf("file skipped due to extension: %s", filePath)
	}
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	if len(content) < MinFileSize {
		return nil, fmt.Errorf("file ignored due to size: %s", filePath)
	}
	return content, nil
}

func (w *Winnowing) GenerateWFP(filePath string) {
	content, err := w.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	fmt.Printf("File: %s, MD5: %s\n", filePath, w.calculateMD5(content))

	// Normalize content and generate k-grams
	normalizedContent := w.normalizeContent(content)
	kgrams := w.generateKgrams(normalizedContent, Gram)

	// Calculate hash for each k-gram
	hashes := make([]uint32, len(kgrams))
	for i, kgram := range kgrams {
		hashes[i] = crc32c(kgram)
	}

	// Select minimum hashes within each window of size `Window`
	var fingerprints []uint32
	for i := 0; i <= len(hashes)-Window; i++ {
		window := hashes[i : i+Window]
		minHash := uint32(MaxCRC32)
		for _, hash := range window {
			if hash < minHash {
				minHash = hash
			}
		}
		fingerprints = append(fingerprints, minHash)
	}

	// Print fingerprints for demonstration
	for i, fingerprint := range fingerprints {
		fmt.Printf("Window %d: Hash %x\n", i, fingerprint)
	}
}

func (w *Winnowing) calculateMD5(content []byte) string {
	hash := md5.Sum(content)
	return hex.EncodeToString(hash[:])
}

// Placeholder for the crc32c function
func crc32c(data []byte) uint32 {
	// This should be replaced with an actual crc32c implementation
	return crc32.ChecksumIEEE(data)
}

func (w *Winnowing) normalizeContent(content []byte) []byte {
	normalized := make([]byte, 0, len(content))
	for _, b := range content {
		if (b >= ASCII0 && b <= ASCII9) || (b >= ASCIILF && b <= ASCIIBackslash) {
			normalized = append(normalized, b)
		}
	}
	return normalized
}

func (w *Winnowing) generateKgrams(content []byte, k int) [][]byte {
	var kgrams [][]byte
	for i := 0; i <= len(content)-k; i++ {
		kgrams = append(kgrams, content[i:i+k])
	}
	return kgrams
}
