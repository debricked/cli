package file

import (
	"strings"
	"testing"
)

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
	if !compiledF.Match("/home/test/regex.test") {
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
	if !compiledF.Match("/home/test/lockFileRegex.test") {
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
	if compiledF.Match("nil") {
		t.Error("failed to assert that no file was matched")
	}
}
