package util

import (
	"errors"
	"strings"
)

// ErrInvalidFormat has a fixed message on purpose. Underlying parser errors can quote file
// contents, and those must not reach the output.
var ErrInvalidFormat = errors.New("content is not valid for this file type")

// ErrNoLockDocument is for a pnpm lock file that holds only pnpm's own environment document.
var ErrNoLockDocument = errors.New("the file has no lock file document, only pnpm's environment document")

// ErrUnsupportedVersion is for a version of a file format that sanitize doesn't handle. The file
// is skipped rather than failed, and no copy of it is written.
var ErrUnsupportedVersion = errors.New("this version of the file format is not supported by sanitize")

// Reason returns the message to report for an error from a FileSanitizer. Only these fixed
// messages can reach the output, never the text of an underlying error.
func Reason(err error) string {
	for _, known := range []error{ErrNoLockDocument, ErrUnsupportedVersion} {
		if errors.Is(err, known) {
			return known.Error()
		}
	}

	return ErrInvalidFormat.Error()
}

// FileSanitizer returns the sanitized content of one file, or ErrInvalidFormat.
type FileSanitizer func(content []byte, context *FileContext) ([]byte, error)

// FileContext gives a file sanitizer the placeholders of the project the file belongs to.
type FileContext struct {
	scope *ProjectScope
	path  string
}

// NewFileContext is for the file at path, relative to the scanned directory.
func NewFileContext(scope *ProjectScope, path string) *FileContext {
	return &FileContext{scope: scope, path: path}
}

func (c *FileContext) ProjectPlaceholder(name string) string {
	return c.scope.projectPlaceholder(name, c.path)
}

func (c *FileContext) FilePlaceholder(name string) string {
	return c.scope.filePlaceholder(name, c.path)
}

func (c *FileContext) PathPlaceholder(format, name string) string {
	return c.scope.pathPlaceholder(format, name, c.path)
}

// RecordPlaceholder is for placeholders that only need to be unique within their file.
func (c *FileContext) RecordPlaceholder(placeholder, name string) {
	c.scope.all.add(placeholder, name, c.path)
}

// SplitLines splits a text file into lines, with Windows line endings removed.
func SplitLines(content []byte) []string {
	return strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")
}
