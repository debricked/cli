package testdata

import (
	"bytes"
	"github.com/debricked/cli/pkg/file"
	"io"
	"net/http"
)

type DebClientMock struct {
	reponseQueue []MockResponse
}

func NewDebClientMock() DebClientMock {
	return DebClientMock{reponseQueue: []MockResponse{}}
}

func (mock *DebClientMock) Post(_ string, _ string, _ *bytes.Buffer) (*http.Response, error) {
	return nil, nil
}

type MockResponse struct {
	StatusCode   int
	ResponseBody io.ReadCloser
	Error        error
}

func (mock *DebClientMock) AddMockResponse(response MockResponse) {
	mock.reponseQueue = append(mock.reponseQueue, response)
}

func (mock *DebClientMock) Get(_ string, _ string) (*http.Response, error) {
	responseMock := mock.reponseQueue[0]      // The first element is the one to be dequeued.
	mock.reponseQueue = mock.reponseQueue[1:] // Slice off the element once it is dequeued.

	res := http.Response{
		Status:           "",
		StatusCode:       responseMock.StatusCode,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             responseMock.ResponseBody,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}

	return &res, responseMock.Error
}

var FormatsMock = []file.Format{
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
