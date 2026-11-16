package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	projectPlaceholderFormat = "debricked-redacted-project-%d"
	filePlaceholderFormat    = "debricked-redacted-file-%d-%d"
)

// MapEntry links a placeholder in the sanitized files to the real name and the file it came from.
type MapEntry struct {
	Placeholder string `json:"placeholder"`
	Name        string `json:"name"`
	Path        string `json:"path"`
}

// Placeholders hands out names that are unique across a whole sanitize run. /sbom/merge matches
// components by purl, so a placeholder shared by two projects would make them look like one.
type Placeholders struct {
	projectCount int
	pathCount    int
	entries      []MapEntry
}

// Entries returns every placeholder handed out so far, in order.
func (p *Placeholders) Entries() []MapEntry {
	return p.entries
}

func (p *Placeholders) add(placeholder, name, path string) string {
	p.entries = append(p.entries, MapEntry{Placeholder: placeholder, Name: name, Path: path})

	return placeholder
}

// ProjectScope gives the files of one project consistent placeholders: the same real name always
// gets the same placeholder within a project, and a different one in any other project.
type ProjectScope struct {
	all       *Placeholders
	number    int
	projects  map[string]string
	files     map[string]string
	fileCount int
}

func NewProjectScope(all *Placeholders, number int) *ProjectScope {
	return &ProjectScope{all: all, number: number, projects: map[string]string{}, files: map[string]string{}}
}

func (s *ProjectScope) projectPlaceholder(name, path string) string {
	if placeholder, found := s.projects[name]; found {
		return placeholder
	}
	s.all.projectCount++
	placeholder := s.all.add(fmt.Sprintf(projectPlaceholderFormat, s.all.projectCount), name, path)
	s.projects[name] = placeholder

	return placeholder
}

// pathPlaceholder numbers across the whole run. In some cases the backend reports a local path as
// a package, so the same placeholder in two projects would merge them.
func (s *ProjectScope) pathPlaceholder(format, name, path string) string {
	s.all.pathCount++

	return s.all.add(fmt.Sprintf(format, s.all.pathCount), name, path)
}

func (s *ProjectScope) filePlaceholder(name, path string) string {
	if placeholder, found := s.files[name]; found {
		return placeholder
	}
	s.fileCount++
	placeholder := s.all.add(fmt.Sprintf(filePlaceholderFormat, s.number, s.fileCount), name, path)
	s.files[name] = placeholder

	return placeholder
}

// WriteMap writes the local map. It holds the real names, so only the current user can read it.
func WriteMap(path string, entries []MapEntry) error {
	if entries == nil {
		entries = []MapEntry{}
	}
	content, err := json.MarshalIndent(entries, "", "    ")
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return err
	}

	return os.WriteFile(path, append(content, '\n'), 0o600)
}
