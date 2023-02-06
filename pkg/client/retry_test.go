package client

import (
	"testing"
	"time"
)

func TestNewRetryClient(t *testing.T) {
	c := NewRetryClient()
	if c.RetryMax != 3 {
		t.Errorf("failed to assert that RetryMax was %d", c.RetryMax)
	}
	if c.RetryWaitMax.Seconds() != (time.Second * 15).Seconds() {
		t.Errorf("failed to assert that RetryWaitMax was %f", (time.Second * 15).Seconds())
	}
	if c.RetryWaitMin.Seconds() != (time.Second * 3).Seconds() {
		t.Errorf("failed to assert that RetryWaitMin was %f", (time.Second * 3).Seconds())
	}
}
