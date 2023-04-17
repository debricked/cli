package file

import (
	"fmt"
	"regexp"
	"regexp/syntax"
	"strings"

	"github.com/debricked/cli/internal/file/pcre"
)

type Format struct {
	ManifestFileRegex string   `json:"regex"`
	DocumentationUrl  string   `json:"documentationUrl"`
	LockFileRegexes   []string `json:"lockFileRegexes"`
}

func NewCompiledFormat(format *Format) (*CompiledFormat, error) {
	var compiledRegex *regexp.Regexp
	var err error
	isPcre := false

	if len(format.ManifestFileRegex) > 0 {
		compiledRegex, err = regexp.Compile(format.ManifestFileRegex)
		if err != nil && strings.Contains(err.Error(), syntax.ErrInvalidPerlOp.String()) {
			isPcre = true
			err = nil
		}
	}

	var lockErr error
	var compiledLockFileRegexes []*regexp.Regexp
	for _, lockFileRegex := range format.LockFileRegexes {
		if len(lockFileRegex) > 0 {
			compiledLockFileRegex, compErr := regexp.Compile(lockFileRegex)
			if compErr == nil {
				compiledLockFileRegexes = append(compiledLockFileRegexes, compiledLockFileRegex)
			} else if strings.Contains(compErr.Error(), syntax.ErrInvalidPerlOp.String()) {
				isPcre = true
			} else {
				lockErr = compErr
			}
		}
	}

	compiledFormat := CompiledFormat{
		compiledRegex,
		&format.DocumentationUrl,
		compiledLockFileRegexes,
		format,
		isPcre,
	}

	if err == nil && lockErr != nil {
		err = lockErr
	}

	return &compiledFormat, err
}

type CompiledFormat struct {
	ManifestFileRegex *regexp.Regexp
	DocumentationUrl  *string
	LockFileRegexes   []*regexp.Regexp
	format            *Format
	pcre              bool
}

func (format *CompiledFormat) MatchFile(filename string) bool {
	if format.pcre {
		matched, err := pcre.Match(format.format.ManifestFileRegex, filename)
		if err != nil {
			fmt.Println(err)
		}

		return matched
	}

	if format.ManifestFileRegex != nil && format.ManifestFileRegex.MatchString(filename) {
		return true
	}

	return false
}

func (format *CompiledFormat) MatchLockFile(filename string) bool {
	if format.pcre {
		for _, lockFileRegex := range format.format.LockFileRegexes {
			matched, _ := pcre.Match(lockFileRegex, filename)
			if matched {
				return true
			}
		}

		return false
	}

	for _, lockFileFormat := range format.LockFileRegexes {
		if lockFileFormat.MatchString(filename) {
			return true
		}
	}

	return false
}
