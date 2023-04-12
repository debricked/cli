package pcre

import (
	"errors"
	"regexp"
)

var (
	NotSupportedErr = errors.New("PCRE regex not supported")
	SyntaxErr       = errors.New("PCRE syntax error")
)

const (
	negLookaheadGroupIdentifierRegex = `\(\?\!`
	negLookaheadGroupRegex           = `^(\(\(\?\!.*\)\))`
	negLookaheadGroupExclusionRegex  = `\(\(\?\!(.+)\)\)`
)

var (
	compiledNegLookaheadGroupIdentifierRegex = regexp.MustCompile(negLookaheadGroupIdentifierRegex)
	compiledNegLookaheadGroupRegex           = regexp.MustCompile(negLookaheadGroupRegex)
	compileNegLookaheadGroupExclusionRegex   = regexp.MustCompile(negLookaheadGroupExclusionRegex)
)

func Match(pcreRegex string, str string) (bool, error) {
	// Evaluate PCRE format
	negLookaheadIds := compiledNegLookaheadGroupIdentifierRegex.FindAllString(pcreRegex, 2)
	if len(negLookaheadIds) != 1 {
		return false, NotSupportedErr
	}
	negLookaheadGroups := compiledNegLookaheadGroupRegex.FindAllString(pcreRegex, 1)
	if len(negLookaheadGroups) != 1 {
		return false, NotSupportedErr
	}
	negLookaheadGroup := negLookaheadGroups[0]

	// Match on negative lookahead
	exclusionRegexes := compileNegLookaheadGroupExclusionRegex.FindStringSubmatch(negLookaheadGroup)
	if len(exclusionRegexes) == 2 {
		matched, err := regexp.MatchString(exclusionRegexes[1], str)
		if err != nil {
			return false, SyntaxErr
		}
		if matched {
			return false, nil
		}
	}

	// Match RE2 compatible regex
	matchingRegex := pcreRegex[len(negLookaheadGroup):]
	compiledMatchingRegex, err := regexp.Compile(matchingRegex)
	if err != nil {
		return false, SyntaxErr
	}

	return compiledMatchingRegex.MatchString(str), nil
}
