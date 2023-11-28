package client

import (
	"bytes"
	"net/http"
	"os"
)

const DefaultDebrickedUri = "https://debricked.com"
const DefaultTimeout = 15

type IDebClient interface {
	// Post makes a POST request to one of Debricked's API endpoints
	Post(uri string, contentType string, body *bytes.Buffer, timeout int) (*http.Response, error)
	// Get makes a GET request to one of Debricked's API endpoints
	Get(uri string, format string) (*http.Response, error)
	SetAccessToken(accessToken *string)
	ConfigureClientSettings(retry bool, timeout int)
}

type DebClient struct {
	host        *string
	httpClient  IClient
	accessToken *string
	jwtToken    string
	retry       bool
	timeout     int
}

func NewDebClient(accessToken *string, httpClient IClient) *DebClient {
	host := os.Getenv("DEBRICKED_URI")
	if len(host) == 0 {
		host = DefaultDebrickedUri
	}

	return &DebClient{
		host:        &host,
		httpClient:  httpClient,
		accessToken: initAccessToken(accessToken),
		jwtToken:    "",
		retry:       true,
		timeout:     DefaultTimeout,
	}
}

func (debClient *DebClient) Post(uri string, contentType string, body *bytes.Buffer, timeout int) (*http.Response, error) {
	if timeout > 0 {
		return postWithTimeout(uri, debClient, contentType, body, true, timeout)
	}

	return post(uri, debClient, contentType, body, true)
}

func (debClient *DebClient) Get(uri string, format string) (*http.Response, error) {
	return get(uri, debClient, debClient.retry, format)
}

func (debClient *DebClient) SetAccessToken(accessToken *string) {
	debClient.accessToken = initAccessToken(accessToken)
}

func (debClient *DebClient) ConfigureClientSettings(retry bool, timeout int) {
	debClient.retry = retry
	debClient.timeout = timeout
}

func initAccessToken(accessToken *string) *string {
	if accessToken == nil {
		accessToken = new(string)
	}
	if len(*accessToken) == 0 {
		*accessToken = os.Getenv("DEBRICKED_TOKEN")
	}

	return accessToken
}
