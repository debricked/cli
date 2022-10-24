package client

import (
	"github.com/hashicorp/go-retryablehttp"
	"time"
)

func newRetryClient() *retryablehttp.Client {
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.RetryWaitMax = time.Second * 15
	client.RetryWaitMin = time.Second * 3
	client.Logger = nil
	return client
}
