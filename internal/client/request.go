package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/fatih/color"
	"github.com/hashicorp/go-retryablehttp"
)

var NoResErr = errors.New("failed to get response. Check out the Debricked status page: https://status.debricked.com/")

func get(uri string, debClient *DebClient, retry bool, format string) (*http.Response, error) {
	request, err := newRequest("GET", *debClient.host+uri, debClient.jwtToken, format, nil)
	if err != nil {
		return nil, err
	}
	res, _ := debClient.httpClient.Do(request)
	req := func() (*http.Response, error) {
		return get(uri, debClient, false, format)
	}

	return interpret(res, req, debClient, retry)
}

func post(uri string, debClient *DebClient, contentType string, body *bytes.Buffer, retry bool) (*http.Response, error) {
	request, err := newRequest("POST", *debClient.host+uri, debClient.jwtToken, "application/json", body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", contentType)
	res, err := debClient.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	req := func() (*http.Response, error) {
		return post(uri, debClient, contentType, body, false)
	}

	return interpret(res, req, debClient, retry)
}

// newRequest creates a new HTTP request with necessary headers added
func newRequest(method string, url string, jwtToken string, format string, body io.Reader) (*retryablehttp.Request, error) {
	req, err := retryablehttp.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", format)
	req.Header.Add("Authorization", "Bearer "+jwtToken)

	return req, nil
}

// interpret a http response
func interpret(res *http.Response, request func() (*http.Response, error), debClient *DebClient, retry bool) (*http.Response, error) {
	if res == nil {
		return nil, NoResErr
	} else if res.StatusCode == http.StatusUnauthorized {
		errMsg := `Unauthorized. Specify access token. 
Read more on https://debricked.com/docs/administration/access-tokens.html`
		if retry {
			err := debClient.authenticate()
			if err != nil {
				return nil, errors.New(errMsg)
			}

			return request()
		}

		return nil, errors.New(errMsg)
	}

	return res, nil
}

type errorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (debClient *DebClient) authenticate() error {
	uri := "/api/login_refresh"

	data := map[string]string{"refresh_token": *debClient.accessToken}
	jsonData, _ := json.Marshal(data)
	res, reqErr := debClient.httpClient.Post(
		*debClient.host+uri,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if reqErr != nil {
		return reqErr
	}

	defer res.Body.Close()
	var tokenData map[string]string
	body, _ := io.ReadAll(res.Body)
	err := json.Unmarshal(body, &tokenData)
	if err != nil {
		var errMessage errorMessage
		_ = json.Unmarshal(body, &errMessage)

		return fmt.Errorf("%s %s\n", color.RedString("⨯"), errMessage.Message)
	}
	debClient.jwtToken = tokenData["token"]

	return nil
}
