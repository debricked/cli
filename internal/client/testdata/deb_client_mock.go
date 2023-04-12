package testdata

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/debricked/cli/internal/client"
)

type DebClientMock struct {
	realDebClient    *client.DebClient
	responseQueue    []MockResponse
	responseUriQueue map[string][]MockResponse
}

func NewDebClientMock() *DebClientMock {
	debClient := client.NewDebClient(nil, nil)

	return &DebClientMock{
		realDebClient:    debClient,
		responseQueue:    []MockResponse{},
		responseUriQueue: map[string][]MockResponse{}}
}

func (mock *DebClientMock) Get(uri string, format string) (*http.Response, error) {
	response, err := mock.popResponse(mock.RemoveQueryParamsFromUri(uri))

	if response != nil {
		return response, err
	}

	return mock.realDebClient.Get(uri, format)
}

func (mock *DebClientMock) Post(uri string, format string, body *bytes.Buffer) (*http.Response, error) {
	response, err := mock.popResponse(mock.RemoveQueryParamsFromUri(uri))

	if response != nil {
		return response, err
	}

	return mock.realDebClient.Post(uri, format, body)
}

type MockResponse struct {
	StatusCode   int
	ResponseBody io.ReadCloser
	Error        error
}

func (mock *DebClientMock) AddMockResponse(response MockResponse) {
	if response.ResponseBody == nil {
		response.ResponseBody = io.NopCloser(bytes.NewReader(nil))
	}
	mock.responseQueue = append(mock.responseQueue, response)
}

func (mock *DebClientMock) AddMockUriResponse(uri string, response MockResponse) {
	_, ok := mock.responseUriQueue[uri]
	if !ok {
		mock.responseUriQueue[uri] = []MockResponse{}
	}

	mock.responseUriQueue[uri] = append(mock.responseUriQueue[uri], response)
}

func (mock *DebClientMock) RemoveQueryParamsFromUri(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return uri
	}
	q := u.Query()
	for s := range q {
		q.Del(s)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

func (mock *DebClientMock) popResponse(uri string) (*http.Response, error) {
	var responseMock MockResponse
	uriQueue, existsInUriQueue := mock.responseUriQueue[uri]
	existsInQueue := len(mock.responseQueue) != 0
	if existsInUriQueue && len(uriQueue) > 0 {
		responseMock = mock.responseUriQueue[uri][0]                // The first element is the one to be dequeued.
		mock.responseUriQueue[uri] = mock.responseUriQueue[uri][1:] // Slice off the element once it is dequeued.
	} else if existsInQueue {
		responseMock = mock.responseQueue[0]        // The first element is the one to be dequeued.
		mock.responseQueue = mock.responseQueue[1:] // Slice off the element once it is dequeued.
	}

	var res http.Response

	if existsInUriQueue || existsInQueue {
		res = http.Response{
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
		}
	} else {
		return nil, errors.New("")
	}

	return &res, responseMock.Error
}
