package file

import (
	"strings"
	"testing"
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
		Regex:            "requirements.*(?:\\.txt)",
		DocumentationUrl: "https://debricked.com/docs/language-support/python.html",
		LockFileRegexes:  nil,
	},
}

func TestNewCompiledFormat(t *testing.T) {
	f := &Format{
		"regex",
		"url",
		[]string{"lockFileRegex"},
	}
	compiledF, err := NewCompiledFormat(f)
	if err != nil {
		t.Error("failed to assert that error was nil")
	}
	if compiledF.Regex.String() != "regex" {
		t.Error("failed to assert that regex was set")
	}
	if *compiledF.DocumentationUrl != "url" {
		t.Error("failed to assert that documentation url was set")
	}
	if len(compiledF.LockFileRegexes) != 1 {
		t.Error("failed to assert that one lock file regex exists")
	}
}
func TestNewCompiledFormatShortRegex(t *testing.T) {
	f := &Format{
		"",
		"",
		[]string{},
	}
	compiledF, err := NewCompiledFormat(f)
	if err == nil {
		t.Error("failed to assert that error occurred")
	}
	if err.Error() != "invalid regex string" {
		t.Error("failed to assert error message")
	}
	if compiledF != nil {
		t.Error("failed to assert that compiled format is nil")
	}
}

func TestNewCompiledFormatBadRegex(t *testing.T) {
	f := &Format{
		")",
		"",
		[]string{},
	}
	compiledF, err := NewCompiledFormat(f)
	if err == nil {
		t.Error("failed to assert that error occurred")
	}
	if strings.Contains("unexpected )", err.Error()) {
		t.Error("failed to assert error message")
	}
	if compiledF != nil {
		t.Error("failed to assert that compiled format is nil")
	}
}

func TestNewCompiledFormatNoLockFileRegex(t *testing.T) {
	f := &Format{
		"regex",
		"url",
		[]string{},
	}
	compiledF, err := NewCompiledFormat(f)
	if err != nil {
		t.Error("failed to assert that error was nil")
	}
	if compiledF.Regex.String() != "regex" {
		t.Error("failed to assert that regex was set")
	}
	if *compiledF.DocumentationUrl != "url" {
		t.Error("failed to assert that documentation url was set")
	}
	if len(compiledF.LockFileRegexes) != 0 {
		t.Error("failed to assert that no lock file regexes exists")
	}
}

func TestNewCompiledFormatBadLockFileRegex(t *testing.T) {
	f := &Format{
		"regex",
		"url",
		[]string{")"},
	}
	compiledF, err := NewCompiledFormat(f)
	if err != nil {
		t.Error("failed to assert that error was nil")
	}
	if compiledF.Regex.String() != "regex" {
		t.Error("failed to assert that regex was set")
	}
	if *compiledF.DocumentationUrl != "url" {
		t.Error("failed to assert that documentation url was set")
	}
	if len(compiledF.LockFileRegexes) != 0 {
		t.Error("failed to assert that no lock file regexes exists")
	}
}

func TestMatchFile(t *testing.T) {
	f := &Format{
		"regex",
		"url",
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
		"regex",
		"url",
		[]string{"lockFileRegex"},
	}
	compiledF, err := NewCompiledFormat(f)
	if err != nil {
		t.Error("failed to assert that error was nil")
	}
	if !compiledF.MatchLockFile("/home/test/lockFileRegex.test") {
		t.Error("failed to find lock file match")
	}
}

func TestMatchNoFile(t *testing.T) {
	f := &Format{
		"regex",
		"url",
		[]string{"lockFileRegex"},
	}
	compiledF, err := NewCompiledFormat(f)
	if err != nil {
		t.Error("failed to assert that error was nil")
	}
	if compiledF.MatchFile("nil") {
		t.Error("failed to assert that no file was matched")
	}
}
