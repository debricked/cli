package auth

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"
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
	defer resp.Body.Close()

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

	go func() {
		awh.Callback(testState)
	}()

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:9096/callback?state=invalid_state&code=test_code")
	if err != nil {
		t.Fatalf("Failed to make callback request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request, got %v", resp.Status)
	}
}

func TestCallbackServerError(t *testing.T) {
	server := &http.Server{Addr: ":9096"}
	go server.ListenAndServe()
	defer server.Shutdown(context.Background())

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
		if err == nil {
			t.Error("Expected an error due to server already running, but got none")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Test timed out")
	}
}

func TestOpenBrowserCmd(t *testing.T) {
	a := NewAuthWebHelper()
	cases := []struct {
		runtimeOS   string
		expectedCmd *exec.Cmd
	}{
		{
			runtimeOS:   "darwin",
			expectedCmd: exec.Command("open", "url"),
		},
		{
			runtimeOS:   "linux",
			expectedCmd: exec.Command("xdg-open", "url"),
		},
		{
			runtimeOS:   "windows",
			expectedCmd: exec.Command("cmd", "/c", "start", "url"),
		},
	}

	for _, c := range cases {
		t.Run(c.runtimeOS, func(t *testing.T) {
			authCmd := a.openBrowserCmd(c.runtimeOS, "url")
			assert.Equal(t, c.expectedCmd.Args, authCmd.Args)
		})
	}
}
