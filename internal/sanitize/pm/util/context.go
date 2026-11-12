package util

import (
	"errors"
	"strings"
)

// ErrInvalidFormat has a fixed message on purpose. Underlying parser errors can quote file
// contents, and those must not reach the output.
var ErrInvalidFormat = errors.New("content is not valid for this file type")

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
