package pcre

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type matchTest struct {
	pcreRegex string
	str       string
	err       error
	matched   bool
}

var matchTests = []matchTest{
	{
		pcreRegex: `((?!WORKSPACE|BUILD)).*(?:\.bazel)`,
		str:       "deps.bazel",
		err:       nil,
		matched:   true,
	},
	{
		pcreRegex: `((?!WORKSPACE|BUILD)).*(?:\.bazel)`,
		str:       "workspace.bzl",
		err:       nil,
		matched:   false,
	},
	{
		pcreRegex: `((?!WORKSPACE|BUILD)).*(?:\.bazel)`,
		str:       "WORKSPACE.bzl",
		err:       nil,
		matched:   false,
	},
	{
		pcreRegex: `((?!WORKSPACE|BUILD)).*(?:\.bzl)`,
		str:       "deps.bzl",
		err:       nil,
		matched:   true,
	},
	{
		pcreRegex: `((?!WORKSPACE|BUILD)).*(?:\.bazel)`,
		str:       "WORKSPACE.bazel",
		err:       nil,
		matched:   false,
	},
	{
		pcreRegex: `((?!WORKSPACE|BUILD)).*(?:\.bazel)`,
		str:       "BUILD.bazel",
		err:       nil,
		matched:   false,
	},
	{
		pcreRegex: `((?!BUILD)).*(?:\.bazel)`,
		str:       "BUILD.bazel",
		err:       nil,
		matched:   false,
	},
	{
		pcreRegex: `((?!BUILD)).*(?:\.bazel)`,
		str:       "BUILD-test.bazel",
		err:       nil,
		matched:   false,
	},
	{
		pcreRegex: `((?>!WORKSPACE|BUILD)).*(?:\.bazel)`,
		str:       "BUILD.bazel",
		err:       NotSupportedErr,
		matched:   false,
	},
	{
		pcreRegex: `((?<!WORKSPACE|BUILD)).*(?:\.bazel)`,
		str:       "BUILD.bazel",
		err:       NotSupportedErr,
		matched:   false,
	},
	{
		pcreRegex: `((?!WORKSPACE|BUILD)).*((?!WORKSPACE|BUILD))(?:\.bazel)`,
		str:       "BUILD.bazel",
		err:       NotSupportedErr,
		matched:   false,
	},
	{
		pcreRegex: `((?WORKSPACE|BUILD)).*((?!WORKSPACE|BUILD))(?:\.bazel)`,
		str:       "BUILD.bazel",
		err:       NotSupportedErr,
		matched:   false,
	},
	{
		pcreRegex: `((?!WORKSPACE|BUILD)).(?pcre-syntax)`,
		str:       "deps.bazel",
		err:       SyntaxErr,
		matched:   false,
	},
	{
		pcreRegex: `((?!WORKSPACE???BUILD)).*(?:\.bazel)`,
		str:       "deps.bazel",
		err:       SyntaxErr,
		matched:   false,
	},
}

func TestMatch(t *testing.T) {
	for _, matchT := range matchTests {
		name := fmt.Sprintf("PCRE ManifestFileRegex: %s, string: %s, matched: %t", matchT.pcreRegex, matchT.str, matchT.matched)
		t.Run(name, func(t *testing.T) {
			matched, err := Match(matchT.pcreRegex, matchT.str)
			assert.Equal(t, matchT.err, err)
			assert.Equal(t, matchT.matched, matched)
		})
	}
}
