package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	testdataClient "github.com/debricked/cli/pkg/client/testdata/client"
)

var client *DebClient

var tkn = "token"

func TestNewDebClientWithTokenParameter(t *testing.T) {
	debClient := NewDebClient(&tkn, nil)
	if *debClient.host != DefaultDebrickedUri {
		t.Error("failed to assert that host was set properly")
	}
	if *debClient.accessToken != tkn {
		t.Error("failed to assert that access token was set properly")
	}
}

func TestNewDebClientWithNilToken(t *testing.T) {
	debClient := NewDebClient(nil, nil)
	if *debClient.host != DefaultDebrickedUri {
		t.Error("failed to assert that host was set properly")
	}
	if debClient.accessToken == nil {
		t.Error("failed to assert that access token was set properly")
	}
}

const debrickedTknEnvVar = "DEBRICKED_TOKEN"

func TestNewDebClientWithTokenEnvVariable(t *testing.T) {
	envVarTkn := "env-tkn"
	oldEnvValue := os.Getenv(debrickedTknEnvVar)
	err := os.Setenv(debrickedTknEnvVar, "env-tkn")
	if err != nil {
		t.Fatalf("failed to set env var %s", debrickedTknEnvVar)
	}
	defer func(key, value string) {
		err := os.Setenv(key, value)
		if err != nil {
			t.Fatalf("failed to reset env var %s", debrickedTknEnvVar)
		}
	}(debrickedTknEnvVar, oldEnvValue)

	accessToken := ""
	debClient := NewDebClient(&accessToken, nil)
	if *debClient.host != DefaultDebrickedUri {
		t.Error("failed to assert that host was set properly")
	}
	if *debClient.accessToken != envVarTkn {
		t.Errorf("failed to assert that access token was set to %s. Got %s", envVarTkn, *debClient.accessToken)
	}
}

func TestNewDebClientWithWithURI(t *testing.T) {
	accessToken := ""
	os.Setenv("DEBRICKED_URI", "https://subdomain.debricked.com")
	debClient := NewDebClient(&accessToken, nil)
	os.Setenv("DEBRICKED_URI", "")
	if *debClient.host != "https://subdomain.debricked.com" {
		t.Error("failed to assert that host was set properly")
	}
}

func TestClientUnauthorized(t *testing.T) {
	clientMock := testdataClient.NewMock()
	clientMock.AddMockResponse(testdataClient.MockResponse{
		StatusCode: http.StatusUnauthorized,
	})
	clientMock.AddMockResponse(testdataClient.MockResponse{
		StatusCode: http.StatusUnauthorized,
	})
	client = NewDebClient(&tkn, clientMock)

	res, err := client.Get("/api/1.0/open/user-profile/is-admin", "application/json")
	if err == nil {
		t.Error("failed to assert client error")
		defer res.Body.Close()
	}

	if !strings.Contains(err.Error(), "Unauthorized. Specify access token") {
		t.Error("Failed to assert unauthorized error message")
	}
}

func TestGet(t *testing.T) {
	clientMock := testdataClient.NewMock()
	clientMock.AddMockResponse(testdataClient.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: io.NopCloser(strings.NewReader(`{"isAdmin": true}`)),
		Error:        nil,
	})
	client = NewDebClient(&tkn, clientMock)

	res, err := client.Get("/api/1.0/open/user-profile/is-admin", "application/json")
	if err != nil {
		t.Fatal("failed to assert that no client error occurred. Error:", err)
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
	clientMock := testdataClient.NewMock()
	clientMock.AddMockResponse(testdataClient.MockResponse{
		StatusCode:   http.StatusForbidden,
		ResponseBody: io.NopCloser(strings.NewReader("{}")),
		Error:        nil,
	})
	client = NewDebClient(&tkn, clientMock)
	data := map[string]bool{"allowSnooze": true}
	jsonData, _ := json.Marshal(data)
	res, err := client.Post(
		"/api/1.0/open/user-permissions/toggle-allow-snooze",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		t.Fatal("failed to assert that no client error occurred. Error:", err)
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusForbidden {
		t.Error("failed to assert that status code was 403")
	}
}

func TestAuthenticate(t *testing.T) {
	tkn = "0501ac404fd1823d0d4c047f957637a912d3b94713ee32a6"
	jwtTkn := "jwt-tkn"
	clientMock := testdataClient.NewMock()
	clientMock.AddMockResponse(testdataClient.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: io.NopCloser(strings.NewReader(fmt.Sprintf(`{"token": "%s"}`, jwtTkn))),
		Error:        nil,
	})
	client = NewDebClient(&tkn, clientMock)
	err := client.authenticate()
	if err != nil {
		t.Fatal("failed to assert that no error occurred")
	}
	if !strings.EqualFold(jwtTkn, client.jwtToken) {
		t.Errorf("failed to assert that the jwt token was properly set to %s. Got %s", jwtTkn, client.jwtToken)
	}
}
