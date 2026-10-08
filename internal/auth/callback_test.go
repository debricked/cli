package auth

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const testState = "test_state"

func TestCallback(t *testing.T) {
	awh := NewAuthWebHelper()

	resultChan := make(chan string)
	go func() {
		result := awh.Callback(testState)
		resultChan <- result
	}()

	time.Sleep(100 * time.Millisecond)

	testCode := "test_code"
	resp, err := http.Get(fmt.Sprintf("http://localhost:9096/callback?state=%s&code=%s", testState, testCode))
	if err != nil {
		t.Fatalf("Failed to make callback request: %v", err)
	}
	defer func() {
		err := resp.Body.Close()
		assert.NoError(t, err)
	}()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	select {
	case result := <-resultChan:
		if result != testCode {
			t.Errorf("Expected code %s, got %s", testCode, result)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
}

func TestCallbackInvalidState(t *testing.T) {
	awh := NewAuthWebHelper()

	done := make(chan struct{})
	go func() {
		awh.Callback(testState)
		close(done)
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:9096/callback?state=invalid_state&code=test_code")
	if err != nil {
		t.Fatalf("Failed to make callback request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}
	accepted, err := http.Get("http://localhost:9096/callback?state=" + testState + "&code=test_code")
	if err != nil {
		t.Fatal(err)
	}
	_ = accepted.Body.Close()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("callback did not release its listener")
	}
}

func TestBrowserLoginCanRepeatOnSameHelper(t *testing.T) {
	helper := NewAuthWebHelper()
	helper.openURL = func(string) error {
		response, err := http.Get("http://127.0.0.1:9096/callback?state=" + testState + "&code=test_code")
		if err != nil {
			return err
		}
		defer func() { assert.NoError(t, response.Body.Close()) }()
		if response.StatusCode != http.StatusOK {
			return fmt.Errorf("unexpected callback status %d", response.StatusCode)
		}

		return nil
	}
	for attempt := 0; attempt < 2; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		code, err := helper.Login(ctx, "https://example.invalid/oauth/authorize", testState)
		cancel()
		assert.NoError(t, err)
		assert.Equal(t, "test_code", code)
	}
}

func TestBrowserLoginCancellationReleasesListener(t *testing.T) {
	helper := NewAuthWebHelper()
	ctx, cancel := context.WithCancel(context.Background())
	helper.openURL = func(string) error {
		cancel()

		return nil
	}
	_, err := helper.Login(ctx, "https://example.invalid/oauth/authorize", testState)
	assert.ErrorIs(t, err, context.Canceled)
	helper.openURL = func(string) error { return fmt.Errorf("browser unavailable") }
	_, err = helper.Login(context.Background(), "https://example.invalid/oauth/authorize", testState)
	assert.ErrorContains(t, err, "browser unavailable")
}

func TestCallbackServerError(t *testing.T) {
	server := &http.Server{Addr: ":9096", ReadHeaderTimeout: time.Second}
	go func() {
		err := server.ListenAndServe()
		assert.Error(t, err) // Two servers trying to run on localhost:9096
	}()
	defer func() {
		err := server.Shutdown(context.Background())
		assert.NoError(t, err)
	}()

	awh := AuthWebHelper{}

	resultChan := make(chan error)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				resultChan <- fmt.Errorf("panic occurred: %v", r)
			} else {
				resultChan <- nil
			}
		}()
		awh.Callback(testState)
	}()

	select {
	case err := <-resultChan:
		assert.Error(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}
}
