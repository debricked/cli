package license

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/debricked/cli/pkg/client/testdata"
	"github.com/go-git/go-git/v5/utils/ioutil"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestOrder(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	addCommitIdMockResponse(&debClientMock)
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusOK})
	reporter := Reporter{DebClient: &debClientMock}
	args := OrderArgs{Email: ""}
	err := reporter.Order(args)
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}
}

func TestOrderBadArgs(t *testing.T) {
	debClientMock := &testdata.DebClientMock{}
	reporter := Reporter{DebClient: debClientMock}
	args := struct{}{}
	err := reporter.Order(args)
	if err != ArgsError {
		t.Error("failed to assert args error")
	}
}

func TestGetCommitError(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	errorAssertion := errors.New("commitError")
	debClientMock.AddMockResponse(testdata.MockResponse{Error: errorAssertion})
	reporter := Reporter{DebClient: &debClientMock}
	args := OrderArgs{Email: ""}
	err := reporter.Order(args)
	if err != errorAssertion {
		t.Error("failed to assert GetCommit error")
	}
}

func TestOrderUnauthorized(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	errorAssertion := errors.New("unauthorized")
	debClientMock.AddMockResponse(testdata.MockResponse{Error: errorAssertion})
	reporter := Reporter{DebClient: &debClientMock}
	args := OrderArgs{Email: ""}
	err := reporter.Order(args)
	if err != errorAssertion {
		t.Error("failed to assert client error")
	}
}

func TestOrderForbidden(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	addCommitIdMockResponse(&debClientMock)
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusForbidden})
	reporter := Reporter{DebClient: &debClientMock}
	args := OrderArgs{Email: ""}
	err := reporter.Order(args)
	if err != SubscriptionError {
		t.Error("failed to assert client error")
	}
}

func TestOrderNotOkResponse(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	addCommitIdMockResponse(&debClientMock)
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusTeapot})
	reporter := Reporter{DebClient: &debClientMock}
	args := OrderArgs{Email: ""}
	err := reporter.Order(args)
	errMsg := fmt.Sprintf("failed to order report. Status code: %d", http.StatusTeapot)
	if !strings.Contains(err.Error(), errMsg) {
		t.Error("failed to assert error message for unknown status code")
	}
}

func TestGetCommitId(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	addCommitIdMockResponse(&debClientMock)
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusOK})
	reporter := Reporter{DebClient: &debClientMock}
	args := OrderArgs{Email: ""}
	err := reporter.Order(args)
	if err != nil {
		t.Error("failed to assert that no error occurred")
	}
}

func TestGetCommitIdUnauthorized(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	errorAssertion := errors.New("unauthorized")
	debClientMock.AddMockResponse(testdata.MockResponse{Error: errorAssertion})
	reporter := Reporter{DebClient: &debClientMock}
	_, err := reporter.getCommitId("")
	if err != errorAssertion {
		t.Error("failed to assert client error")
	}
}

func TestGetCommitIdForbidden(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusForbidden})
	reporter := Reporter{DebClient: &debClientMock}
	_, err := reporter.getCommitId("")
	if err != SubscriptionError {
		t.Error("failed to assert subscription error")
	}
}

func TestGetCommitIdNotOkResponse(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusTeapot})
	reporter := Reporter{DebClient: &debClientMock}
	_, err := reporter.getCommitId("")
	if !strings.Contains(err.Error(), "No commit was found with the name") {
		t.Error("failed to assert that not commit error message")
	}
}

func TestGetCommitIdNoResult(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: createIoReadCloserFromCommit(nil),
	})
	reporter := Reporter{DebClient: &debClientMock}
	_, err := reporter.getCommitId("")
	if !strings.Contains(err.Error(), "No commit was found with the name") {
		t.Error("failed to assert that not commit error message")
	}
}

func addCommitIdMockResponse(mockClient *testdata.DebClientMock) {
	c := commit{
		FileIds:     []int{},
		Id:          0,
		Name:        "commit-hash",
		ReleaseData: "",
	}
	mockResponse := testdata.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: createIoReadCloserFromCommit(&c),
		Error:        nil,
	}
	mockClient.AddMockResponse(mockResponse)
}

func createIoReadCloserFromCommit(c *commit) io.ReadCloser {
	var commitResponse []commit
	if c != nil {
		commitResponse = append(commitResponse, *c)
	}
	body, _ := json.Marshal(commitResponse)
	reader := bytes.NewReader(body)
	return ioutil.NewReadCloser(reader, nil)
}
