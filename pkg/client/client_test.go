package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

var client *DebClient

func setUp(authorized bool) {
	var token string
	if authorized {
		token = os.Getenv("DEBRICKED_TOKEN")
	} else {
		token = "invalid"
	}
	client = NewDebClient(&token)
}

func TestNewDebClientWithTokenParameter(t *testing.T) {
	accessToken := "token"

	debClient := NewDebClient(&accessToken)
	if debClient == nil {
		t.Error("failed to assert that debricked client was not nil")
	}
	if *debClient.host != "https://debricked.com" {
		t.Error("failed to assert that host was set properly")
	}
	if *debClient.accessToken != accessToken {
		t.Error("failed to assert that access token was set properly")
	}
}

func TestNewDebClientWithNilToken(t *testing.T) {
	debClient := NewDebClient(nil)
	if debClient == nil {
		t.Error("failed to assert that debricked client was not nil")
	}
	if *debClient.host != "https://debricked.com" {
		t.Error("failed to assert that host was set properly")
	}
	if debClient.accessToken == nil {
		t.Error("failed to assert that access token was set properly")
	}
}

func TestNewDebClientWithTokenEnvVariable(t *testing.T) {
	accessToken := ""
	debClient := NewDebClient(&accessToken)
	if debClient == nil {
		t.Error("failed to assert that debricked client was not nil")
	}
	if *debClient.host != "https://debricked.com" {
		t.Error("failed to assert that host was set properly")
	}
	if len(*debClient.accessToken) == 0 {
		t.Error("failed to assert that access token was set properly")
	}
}

func TestNewDebClientWithWithURI(t *testing.T) {
	accessToken := ""
	os.Setenv("DEBRICKED_URI", "https://subdomain.debricked.com")
	debClient := NewDebClient(&accessToken)
	os.Setenv("DEBRICKED_URI", "")

	if debClient == nil {
		t.Error("failed to assert that debricked client was not nil")
	}
	if *debClient.host != "https://subdomain.debricked.com" {
		t.Error("failed to assert that host was set properly")
	}
}

func TestClientUnauthorized(t *testing.T) {
	setUp(false)
	_, err := client.Get("/api/1.0/open/user-profile/is-admin", "application/json")
	if err == nil {
		t.Error("failed to assert client error")
	}
	if !strings.Contains(err.Error(), "Unauthorized. Specify access token") {
		t.Error("Failed to assert unauthorized error message")
	}
}

func TestGet(t *testing.T) {
	setUp(true)
	res, err := client.Get("/api/1.0/open/user-profile/is-admin", "application/json")
	if err != nil {
		t.Fatal("failed to assert that no client error occurred. Error: ", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		t.Error("failed to assert that status code was 200")
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error("failed to read body")
	}
	data := string(body)
	if !strings.Contains(data, "isAdmin") {
		t.Error("failed to assert data contained isAdmin")
	}
}

func TestPost(t *testing.T) {
	setUp(true)
	data := map[string]bool{"allowSnooze": true}
	jsonData, _ := json.Marshal(data)
	res, err := client.Post(
		"/api/1.0/open/user-permissions/toggle-allow-snooze",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		t.Fatal("failed to assert that no client error occurred. Error: ", err.Error())
	}
	if res.StatusCode != http.StatusForbidden {
		t.Error("failed to assert that status code was 403")
	}
}

func TestAuthenticate(t *testing.T) {
	token := "0501ac404fd1823d0d4c047f957637a912d3b94713ee32a6"
	client = NewDebClient(&token)
	err := client.authenticate()
	if err == nil {
		t.Fatal("failed to assert that error occurred")
	}

	if !strings.Contains(err.Error(), "An authentication exception occurred") {
		t.Error("failed to assert that an authentication occurred")
	}
}
