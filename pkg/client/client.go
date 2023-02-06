package client

import (
	"github.com/hashicorp/go-retryablehttp"
	"net/http"
)

type IClient interface {
	Do(req *retryablehttp.Request) (*http.Response, error)
	Post(url, bodyType string, body interface{}) (*http.Response, error)
}
