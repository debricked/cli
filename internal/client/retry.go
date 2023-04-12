package client

import (
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func NewRetryClient() *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.RetryWaitMax = time.Second * 15
	client.RetryWaitMin = time.Second * 3
	client.Logger = nil

	return client
}
