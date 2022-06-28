package client

import (
	"bytes"
	"net/http"
	"os"
)

type Client interface {
	// Post makes a POST request to one of Debricked's API endpoints
	Post(uri string, contentType string, body *bytes.Buffer) (*http.Response, error)
	// Get makes a GET request to one of Debricked's API endpoints
	Get(uri string, format string) (*http.Response, error)
}

type DebClient struct {
	host        *string
	httpClient  *http.Client
	accessToken *string
	jwtToken    string
}

func NewDebClient(accessToken *string) *DebClient {
	if accessToken == nil {
		accessToken = new(string)
	}
	if len(*accessToken) == 0 {
		*accessToken = os.Getenv("DEBRICKED_TOKEN")
	}
	host := os.Getenv("DEBRICKED_URI")
	if len(host) == 0 {
		host = "https://debricked.com"
	}

	return &DebClient{
		host:        &host,
		httpClient:  &http.Client{},
		accessToken: accessToken,
		jwtToken:    "",
	}
}

func (debClient *DebClient) Post(uri string, contentType string, body *bytes.Buffer) (*http.Response, error) {
	return post(uri, debClient, contentType, body, true)
}

func (debClient *DebClient) Get(uri string, format string) (*http.Response, error) {
	return get(uri, debClient, true, format)
}
