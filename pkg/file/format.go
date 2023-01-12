package file

import (
	"regexp"
)

type Format struct {
	Regex            string   `json:"regex"`
	DocumentationUrl string   `json:"documentationUrl"`
	LockFileRegexes  []string `json:"lockFileRegexes"`
}

func NewCompiledFormat(format *Format) (*CompiledFormat, error) {
	var compiledRegex *regexp.Regexp
	var err error
	if len(format.Regex) > 0 {
		compiledRegex, err = regexp.Compile(format.Regex)
	}
	var compiledLockFileRegexes []*regexp.Regexp
	for _, lockFileRegex := range format.LockFileRegexes {
		if len(lockFileRegex) > 0 {
			compiledLockFileRegex, err := regexp.Compile(lockFileRegex)
			if err == nil {
				compiledLockFileRegexes = append(compiledLockFileRegexes, compiledLockFileRegex)
			}
		}
	}

	compiledFormat := CompiledFormat{
		compiledRegex,
		&format.DocumentationUrl,
		compiledLockFileRegexes,
	}

	return &compiledFormat, err
}

type CompiledFormat struct {
	Regex            *regexp.Regexp
	DocumentationUrl *string
	LockFileRegexes  []*regexp.Regexp
}

func (format *CompiledFormat) MatchFile(filename string) bool {
	if format.Regex != nil && format.Regex.MatchString(filename) {
		return true
	}

	return false
}

func (format *CompiledFormat) MatchLockFile(filename string) bool {
	for _, lockFileFormat := range format.LockFileRegexes {
		if lockFileFormat.MatchString(filename) {
			return true
		}
	}

	return false
}
