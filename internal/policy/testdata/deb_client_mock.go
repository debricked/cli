package testdata

import (
	"bytes"
	"net/http"

	clientTestdata "github.com/debricked/cli/internal/client/testdata"
)

// RecordingDebClientMock records the bodies posted to Debricked so that tests
// can assert on the request payload.
type RecordingDebClientMock struct {
	*clientTestdata.DebClientMock
	Uri  string
	Body string
}

func NewRecordingDebClientMock() *RecordingDebClientMock {
	return &RecordingDebClientMock{DebClientMock: clientTestdata.NewDebClientMock()}
}

func (mock *RecordingDebClientMock) Post(uri string, format string, body *bytes.Buffer, timeout int) (*http.Response, error) {
	mock.Uri = uri
	if body != nil {
		mock.Body = body.String()
	}

	return mock.DebClientMock.Post(uri, format, body, timeout)
}
