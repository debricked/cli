package find

import (
	"bytes"
	"debricked/pkg/client"
	"debricked/pkg/file"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

var clientMock client.Client = &debClientMock{}

func TestNewFindCmd(t *testing.T) {
	clientMockAuthorized = true
	cmd := NewFindCmd(&clientMock)

	commands := cmd.Commands()
	nbrOfCommands := 0
	if len(commands) != nbrOfCommands {
		t.Error(fmt.Sprintf("failed to assert that there were %d sub commands connected", nbrOfCommands))
	}

	flags := cmd.Flags()
	flagAssertions := map[string]string{}
	for name, shorthand := range flagAssertions {
		flag := flags.Lookup(name)
		if flag == nil {
			t.Error(fmt.Sprintf("failed to assert that %s flag was set", name))
		}
		if flag.Shorthand != shorthand {
			t.Error(fmt.Sprintf("failed to assert that %s flag shorthand %s was set correctly", name, shorthand))
		}
	}
}

func TestRun(t *testing.T) {
	clientMockAuthorized = true
	finder, _ = file.NewFinder(clientMock)
	err := run(nil, []string{"."})
	if err != nil {
		t.Fatal("failed to assert that no error occurred. Error:", err)
	}
}

func TestFind(t *testing.T) {
	clientMockAuthorized = true
	finder, _ = file.NewFinder(clientMock)
	err := find("../../scan", []string{})
	if err != nil {
		t.Fatal("failed to assert that no error occurred. Error:", err)
	}
}

func TestValidateArgs(t *testing.T) {
	err := validateArgs(nil, []string{"."})
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}
}

func TestValidateArgsInvalidArgs(t *testing.T) {
	err := validateArgs(nil, []string{})
	if err == nil {
		t.Error("failed to assert that an error occurred")
	}
	if !strings.Contains(err.Error(), "requires path") {
		t.Error("failed to assert error message")
	}

	err = validateArgs(nil, []string{"invalid-path"})
	if err == nil {
		t.Error("failed to assert that an error occurred")
	}
	if !strings.Contains(err.Error(), "invalid path specified") {
		t.Error("failed to assert error message")
	}
}

type debClientMock struct{}

func (mock *debClientMock) Post(_ string, _ string, _ *bytes.Buffer) (*http.Response, error) {
	return nil, nil
}

var clientMockAuthorized bool

func (mock *debClientMock) Get(_ string, _ string) (*http.Response, error) {
	var statusCode int
	var body io.ReadCloser = nil
	if clientMockAuthorized {
		statusCode = http.StatusOK
		formatsBytes, _ := json.Marshal(formatsMock)
		body = ioutil.NopCloser(strings.NewReader(string(formatsBytes)))
	} else {
		statusCode = http.StatusForbidden
	}
	res := http.Response{
		Status:           "",
		StatusCode:       statusCode,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             body,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}

	return &res, nil
}

var formatsMock = []file.Format{
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
