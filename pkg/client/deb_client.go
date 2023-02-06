package client

import (
	"bytes"
	"net/http"
	"os"
)

const DefaultDebrickedUri = "https://debricked.com"

type IDebClient interface {
	// Post makes a POST request to one of Debricked's API endpoints
	Post(uri string, contentType string, body *bytes.Buffer) (*http.Response, error)
	// Get makes a GET request to one of Debricked's API endpoints
	Get(uri string, format string) (*http.Response, error)
}

type DebClient struct {
	host        *string
	httpClient  IClient
	accessToken *string
	jwtToken    string
}

func NewDebClient(accessToken *string, httpClient IClient) *DebClient {
	if accessToken == nil {
		accessToken = new(string)
	}
	if len(*accessToken) == 0 {
		*accessToken = os.Getenv("DEBRICKED_TOKEN")
	}
	host := os.Getenv("DEBRICKED_URI")
	if len(host) == 0 {
		host = DefaultDebrickedUri
	}

	return &DebClient{
		host:        &host,
		httpClient:  httpClient,
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
