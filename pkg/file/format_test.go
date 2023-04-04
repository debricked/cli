package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	url   = "url"
	regex = "regex"
)

var formatsMock = []Format{
	{
		// Format with regex and lock file regex
		Regex:            "composer\\.json",
		DocumentationUrl: "https://debricked.com/docs/language-support/php.html",
		LockFileRegexes:  []string{"composer\\.lock"},
	},
	{
		// Format with regex and multiple lock file regexes
		Regex:            "package\\.json",
		DocumentationUrl: "https://debricked.com/docs/language-support/javascript.html",
		LockFileRegexes:  []string{"yarn\\.lock", "package-lock\\.json"},
	},
	{
		// Format with regex and debricked made lock file regex
		Regex:            "go\\.mod",
		DocumentationUrl: "https://debricked.com/docs/language-support/golang.html",
		LockFileRegexes:  []string{"\\.debricked-go-dependencies\\.txt"},
	},
	{
		// Format without regex but with one lock file regex
		Regex:            "",
		DocumentationUrl: "https://debricked.com/docs/language-support/rust.html",
		LockFileRegexes:  []string{"Cargo\\.lock"},
	},
	{
		// Format with regex but without lock file regexes
		Regex:            "requirements.*\\.txt$",
		DocumentationUrl: "https://debricked.com/docs/language-support/python.html",
		LockFileRegexes:  []string{".*\\.pip\\.debricked\\.lock"},
	},
}

func TestNewCompiledFormat(t *testing.T) {
	f := &Format{
		regex,
		url,
		[]string{"lockFileRegex"},
	}

	compiledF, err := NewCompiledFormat(f)

	assert.NoError(t, err)
	assert.Equal(t, regex, compiledF.Regex.String())
	assert.Equal(t, url, *compiledF.DocumentationUrl, "failed to assert that documentation url was set")
	assert.Len(t, compiledF.LockFileRegexes, 1, "failed to assert that one lock file regex exists")
}

func TestNewCompiledFormatNoRegexes(t *testing.T) {
	f := &Format{
		"",
		"",
		[]string{},
	}

	compiledF, err := NewCompiledFormat(f)

	assert.NoError(t, err)
	assert.Nil(t, compiledF.Regex)
	assert.Len(t, compiledF.LockFileRegexes, 0)
}

func TestNewCompiledFormatBadRegex(t *testing.T) {
	f := &Format{
		")",
		"",
		[]string{},
	}

	compiledF, err := NewCompiledFormat(f)

	assert.Error(t, err)
	assert.ErrorContains(t, err, "unexpected )")
	assert.NotNil(t, compiledF)
}

func TestNewCompiledFormatNoLockFileRegex(t *testing.T) {
	f := &Format{
		regex,
		url,
		[]string{},
	}
	compiledF, err := NewCompiledFormat(f)
	if err != nil {
		t.Error("failed to assert that error was nil")
	}
	if compiledF.Regex.String() != regex {
		t.Error("failed to assert that regex was set")
	}
	if *compiledF.DocumentationUrl != url {
		t.Error("failed to assert that documentation url was set")
	}
	if len(compiledF.LockFileRegexes) != 0 {
		t.Error("failed to assert that no lock file regexes exists")
	}
}

func TestNewCompiledFormatBadLockFileRegex(t *testing.T) {
	f := &Format{
		regex,
		url,
		[]string{")"},
	}
	compiledF, err := NewCompiledFormat(f)
	assert.Error(t, err)
	assert.ErrorContains(t, err, "unexpected )")
	assert.NotNil(t, compiledF)
}

func TestNewCompiledFormatNoFileRegex(t *testing.T) {
	f := &Format{
		"",
		url,
		[]string{regex},
	}
	compiledF, err := NewCompiledFormat(f)
	assert.Nil(t, err)
	assert.Nil(t, compiledF.Regex)
	assert.Equal(t, url, *compiledF.DocumentationUrl)
	assert.Len(t, compiledF.LockFileRegexes, 1, "failed to assert that one lock file regex exists")
}

func TestNewCompiledFormatPcre(t *testing.T) {
	f := &Format{
		"(?!.+)",
		"",
		[]string{"(?!.+)"},
	}
	compiledF, err := NewCompiledFormat(f)
	assert.NoError(t, err)
	assert.True(t, compiledF.pcre, "failed to assert that the pcre was set to true")
}

func TestMatchFile(t *testing.T) {
	f := &Format{
		regex,
		url,
		[]string{},
	}
	compiledF, err := NewCompiledFormat(f)
	if err != nil {
		t.Error("failed to assert that error was nil")
	}
	if !compiledF.MatchFile("/home/test/regex.test") {
		t.Error("failed to find filename match")
	}
}

func TestMatchLockFile(t *testing.T) {
	f := &Format{
		regex,
		url,
		[]string{"lockFileRegex"},
	}
	compiledF, err := NewCompiledFormat(f)
	assert.NoError(t, err)
	assert.True(t, compiledF.MatchLockFile("/home/test/lockFileRegex.test"))
}

func TestMatchNoFile(t *testing.T) {
	f := &Format{
		regex,
		url,
		[]string{"lockFileRegex"},
	}
	compiledF, err := NewCompiledFormat(f)
	assert.NoError(t, err)
	assert.False(t, compiledF.MatchFile("nil"))
}

func TestMatchPcreFile(t *testing.T) {
	f := &Format{
		`((?!WORKSPACE|BUILD)).*(?:\.bazel)`,
		url,
		[]string{},
	}
	compiledF, err := NewCompiledFormat(f)
	assert.NoError(t, err)
	assert.True(t, compiledF.MatchFile("deps.bazel"))
}

func TestMatchPcreLockFile(t *testing.T) {
	f := &Format{
		"",
		url,
		[]string{`((?!WORKSPACE|BUILD)).*(?:\.bzl)`},
	}
	compiledF, err := NewCompiledFormat(f)
	assert.NoError(t, err)
	assert.True(t, compiledF.MatchLockFile("deps.bzl"))
}
