package testdata

import (
	"bytes"
	"io"
	"net/http"
)

type DebClientMock struct {
	responseQueue []MockResponse
}

func NewDebClientMock() *DebClientMock {
	return &DebClientMock{responseQueue: []MockResponse{}}
}

func (mock *DebClientMock) Get(_ string, _ string) (*http.Response, error) {
	return mock.popResponse()
}

func (mock *DebClientMock) Post(_ string, _ string, _ *bytes.Buffer) (*http.Response, error) {
	return mock.popResponse()
}

type MockResponse struct {
	StatusCode   int
	ResponseBody io.ReadCloser
	Error        error
}

func (mock *DebClientMock) AddMockResponse(response MockResponse) {
	mock.responseQueue = append(mock.responseQueue, response)
}

func (mock *DebClientMock) popResponse() (*http.Response, error) {
	responseMock := mock.responseQueue[0]       // The first element is the one to be dequeued.
	mock.responseQueue = mock.responseQueue[1:] // Slice off the element once it is dequeued.

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
