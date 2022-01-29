package file

import (
	"errors"
	"regexp"
)

type Format struct {
	Regex            string   `json:"regex"`
	DocumentationUrl string   `json:"documentationUrl"`
	LockFileRegexes  []string `json:"lockFileRegexes"`
}

func NewCompiledFormat(format *Format) (*CompiledFormat, error) {
	if len(format.Regex) < 1 {
		return nil, errors.New("invalid regex string")
	}

	compiledRegex, err := regexp.Compile(format.Regex)
	if err != nil {
		return nil, err
	}

	var compiledLockFileRegexes []*regexp.Regexp
	for _, lockFileRegex := range format.LockFileRegexes {
		if len(format.Regex) > 0 {
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

	return &compiledFormat, nil
}

type CompiledFormat struct {
	Regex            *regexp.Regexp
	DocumentationUrl *string
	LockFileRegexes  []*regexp.Regexp
}

func (format *CompiledFormat) Match(filename string) bool {
	if format.Regex.MatchString(filename) {
		return true
	}

	for _, lockFileFormat := range format.LockFileRegexes {
		if lockFileFormat.MatchString(filename) {
			return true
		}
	}

	return false
}
