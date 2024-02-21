package fingerprint

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"path/filepath"
)

const (
	Gram               = 30
	Window             = 64
	MinFileSizeDefault = 256
	MaxLongLineChars   = 1000
	ASCII0             = 48
	ASCII9             = 57
	ASCIIA             = 65
	ASCIIZ             = 90
	ASCIIa             = 97
	ASCIIz             = 122
	ASCIILF            = 10
	ASCIIBackslash     = 92
	MaxCRC32           = uint32(4294967295)
)

var IncludedExtensions = map[string]bool{
	// C
	".c": true,
	".h": true,

	// C++
	".cc":  true,
	".cpp": true,
	".hpp": true,

	// C#
	".cs": true,

	// Go
	".go": true,

	// Java
	".java": true,

	// Kotlin
	".kt": true,

	// JavaScript + TypeScript + frameworks
	".js":        true,
	".ts":        true,
	".jsx":       true,
	".tsx":       true,
	".vue":       true,
	".svelte":    true,
	".elm":       true,
	".coffee":    true,
	".litcoffee": true,
	".cjsx":      true,
	".iced":      true,
	".es":        true,
	".es6":       true,
	".mjs":       true,

	// Ruby
	".rb": true,

	// Rust
	".rs": true,

	// Swift
	".swift": true,

	// Objective-C
	".m":  true,
	".mm": true,

	// PHP
	".php": true,

	// Python
	".py": true,

	// CSS
	".css": true,

	// Scala
	".scala": true,
}

type Winnowing struct {
	crc8MaximTable []uint8
	MinFileSize    int
	results        *[]Snippet
}

type Snippet struct {
	Content string
	Hash    string
	Line    int
}

func NewWinnowing(minFileSize *int) *Winnowing {
	var MinFileSize int
	if minFileSize != nil {
		MinFileSize = *minFileSize
	} else {
		MinFileSize = MinFileSizeDefault
	}
	return &Winnowing{
		crc8MaximTable: make([]uint8, 0),
		MinFileSize:    MinFileSize,
	}
}

func (w *Winnowing) NormalizeByte(b byte) byte {
	if b < ASCII0 || b > ASCII9 {
		return 0
	}
	return b
}

func (w *Winnowing) ShouldSkipFile(filePath string) bool {
	ext := filepath.Ext(filePath)
	if _, ok := IncludedExtensions[ext]; !ok {
		return true
	}
	return false
}

func (w *Winnowing) Write(p []byte) (n int, err error) {
	var output []Snippet

	content := p
	content_len := len(content)
	if content_len < w.MinFileSize {
		return len(p), nil
	}
	line := 1
	window := []uint32{}
	gram := []byte{}
	last_hash := MaxCRC32
	last_content_window_end := 0
	for i, bt := range content {

		if bt == ASCIILF {
			line++
			continue
		}

		btNorm, process := w.normalizeContent(bt)
		if !process {
			continue
		}

		gram = append(gram, btNorm)

		if len(gram) >= Gram {
			gramCrc32 := crc32c(gram)
			window = append(window, gramCrc32)

			if len(window) >= Window {
				minHash := minHash(window)

				if minHash != last_hash {

					// Hashing the hash to balance the distribution
					crc := crc32c([]byte{byte(minHash & 0xff), byte((minHash >> 8) & 0xff), byte((minHash >> 16) & 0xff), byte((minHash >> 24) & 0xff)})
					output = append(output, Snippet{Content: string(content[last_content_window_end:i]), Hash: fmt.Sprintf("%x", crc), Line: line})
					last_content_window_end = i
				}
				last_hash = minHash
				window = window[1:]
			}
			gram = gram[1:]
		}

	}

	w.results = &output

	return len(p), nil
}

func minHash(window []uint32) uint32 {
	min := MaxCRC32
	for _, hash := range window {
		if hash < min {
			min = hash
		}
	}
	return min
}

func (w *Winnowing) GenerateWFP(filePath string) (*[]Snippet, error) {
	rc, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	if w.ShouldSkipFile(filePath) {
		return nil, nil
	}

	_, err = io.Copy(w, rc)
	if err != nil {
		return nil, err
	}

	return w.results, nil
}

func crc32c(data []byte) uint32 {
	// Create a table for the Castagnoli polynomial.
	castagnoliTable := crc32.MakeTable(crc32.Castagnoli)

	// crc32.ChecksumIEEE(data)
	return crc32.Checksum(data, castagnoliTable)
}

func (w *Winnowing) normalizeContent(b byte) (byte, bool) {
	if b < ASCII0 || b > ASCIIz {
		return 0, false
	} else if b <= ASCII9 || b >= ASCIIa {
		return b, true
	} else if b >= 65 && b <= 90 {
		return b + 32, true
	} else {
		return 0, false
	}
}

func (w *Winnowing) generateKgrams(content []byte, k int) [][]byte {
	var kgrams [][]byte
	for i := 0; i <= len(content)-k; i++ {
		kgrams = append(kgrams, content[i:i+k])
	}
	return kgrams
}
