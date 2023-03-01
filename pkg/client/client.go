package client

import (
	"net/http"

	"github.com/hashicorp/go-retryablehttp"
)

type IClient interface {
	Do(req *retryablehttp.Request) (*http.Response, error)
	Post(url, bodyType string, body interface{}) (*http.Response, error)
}
