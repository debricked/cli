package testdata

import (
	"bytes"
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
