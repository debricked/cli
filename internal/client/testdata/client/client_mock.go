package client

import (
	"bytes"
	"io"
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type Mock struct {
	responseQueue []MockResponse
}

func NewMock() *Mock {
	return &Mock{[]MockResponse{}}
}

func (mock *Mock) Do(_ *retryablehttp.Request) (*http.Response, error) {
	return mock.popResponse()
}

func (mock *Mock) Post(_, _ string, _ interface{}) (*http.Response, error) {
	return mock.popResponse()
}

type MockResponse struct {
	StatusCode   int
	ResponseBody io.ReadCloser
	Error        error
}

func (mock *Mock) AddMockResponse(response MockResponse) {
	if response.ResponseBody == nil {
		response.ResponseBody = io.NopCloser(bytes.NewReader(nil))
	}
	mock.responseQueue = append(mock.responseQueue, response)
}

func (mock *Mock) popResponse() (*http.Response, error) {
	var responseMock MockResponse
	if len(mock.responseQueue) > 0 {
		responseMock = mock.responseQueue[0]        // The first element is the one to be dequeued.
		mock.responseQueue = mock.responseQueue[1:] // Slice off the element once it is dequeued.

		return &http.Response{
			Status:           "",
			StatusCode:       responseMock.StatusCode,
			Proto:            "",
			ProtoMajor:       0,
			ProtoMinor:       0,
			Header:           nil,
			Body:             responseMock.ResponseBody,
			ContentLength:    0,
			TransferEncoding: nil,
			Close:            true,
			Uncompressed:     false,
			Trailer:          nil,
			Request:          nil,
			TLS:              nil,
		}, responseMock.Error
	}

	return nil, nil
}
